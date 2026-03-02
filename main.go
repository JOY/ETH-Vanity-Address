package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/term"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/planxnx/ethereum-wallet-generator/internal/encryption"
	enginefactory "github.com/planxnx/ethereum-wallet-generator/internal/engine"
	"github.com/planxnx/ethereum-wallet-generator/internal/generators"
	"github.com/planxnx/ethereum-wallet-generator/internal/progressbar"
	"github.com/planxnx/ethereum-wallet-generator/internal/repository"
	"github.com/planxnx/ethereum-wallet-generator/utils"
	"github.com/planxnx/ethereum-wallet-generator/wallets"
)

func init() {
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		if err := os.Mkdir("db", 0o750); err != nil {
			panic(err)
		}
	}
}

func main() {
	// Context with gracefully shutdown signal
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)
	defer stop()

	fmt.Println("===============ETH Wallet Generator===============")
	fmt.Println(" ")

	// Parse flags
	decryptFile := flag.String("decrypt", "", "decrypt encrypted JSON file (e.g., wallets.encrypted.json)")
	number := flag.Int("n", 10, "set number of generate times (not number of result wallets) (set number to -1 for Infinite loop ∞)")
	limit := flag.Int("limit", 0, "set limit number of result wallets. stop generate when result of vanity wallets reach the limit (set number to 0 for no limit)")
	dbPath := flag.String("db", "", "set sqlite output name eg. wallets.db (db file will create in /db)")
	concurrency := flag.Int("c", 1, "set concurrency value")
	bits := flag.Int("bit", 128, "set number of entropy bits [128, 256]")
	strict := flag.Bool("strict", false, "strict contains mode (required contains to use)")
	contain := flag.String("contains", "", "show only result that contained with the given letters (support for multiple characters)")
	prefix := flag.String("prefix", "", "show only result that prefix was matched with the given letters  (support for single character)")
	suffix := flag.String("suffix", "", "show only result that suffix was matched with the given letters (support for single character)")
	regEx := flag.String("regex", "", "show only result that was matched with given regex (eg. ^0x99 or ^0x00)")
	isDryrun := flag.Bool("dryrun", false, "generate wallet without a result (used for benchmark speed)")
	isCompatible := flag.Bool("compatible", false, "logging compatible mode (turn this on to fix logging glitch)")
	mode := flag.Int("mode", 1, "wallet generate mode [1: normal mode, 2: only private key mode(generate only privatekey, this fastest mode)]")
	engineName := flag.String("engine", "cpu", "wallet generation engine [cpu|gpu]")
	gpuBin := flag.String("gpu-bin", "", "path to external GPU worker binary (optional; defaults to vanity-eth-address)")
	gpuArgs := flag.String("gpu-args", "", "arguments for GPU worker binary, space-delimited (used when -engine gpu)")
	outputMode := flag.String("output-mode", "safe", "terminal output mode [safe|full|silent]")
	outputDir := flag.String("output-dir", "output", "directory for encrypted export file (in-memory mode)")
	outputName := flag.String("output-name", "", "encrypted export filename; supports {timestamp} token (in-memory mode)")
	noExport := flag.Bool("no-export", false, "disable encrypted export on exit (in-memory mode)")
	contractMode := flag.String("contract-mode", "", "gpu contract search mode [create|create2|create3] (optional)")
	contractBytecode := flag.String("contract-bytecode", "", "path to contract bytecode file (required for create2/create3)")
	contractAddress := flag.String("contract-address", "", "origin/sender address (required for create2/create3)")
	contractDeployer := flag.String("contract-deployer", "", "deployer address (required for create3)")
	flag.Parse()

	// Handle decrypt command
	if *decryptFile != "" {
		// Read encrypted file
		encryptedData, err := os.ReadFile(*decryptFile)
		if err != nil {
			log.Fatalf("Error reading encrypted file: %v", err)
		}

		// Prompt for password
		fmt.Print("Enter password to decrypt: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println() // newline after password input
		if err != nil {
			log.Fatalf("Error reading password: %v", err)
		}
		if len(password) == 0 {
			log.Fatal("Error: password cannot be empty")
		}

		// Parse JSON (file is not fully encrypted, only fields are encrypted)
		var encryptedWallets []repository.EncryptedWallet
		if err := json.Unmarshal(encryptedData, &encryptedWallets); err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		// Decrypt individual fields and convert to Wallet
		decryptedWallets := make([]*wallets.Wallet, 0, len(encryptedWallets))
		for _, encWallet := range encryptedWallets {
			wallet := &wallets.Wallet{
				Address: encWallet.Address,
				HDPath:  encWallet.HDPath,
				Bits:    encWallet.Bits,
			}
			wallet.CreatedAt = encWallet.CreatedAt
			wallet.UpdatedAt = encWallet.UpdatedAt

			// Decrypt PrivateKey if present
			if encWallet.PrivateKey != "" {
				decryptedPK, err := encryption.Decrypt(encWallet.PrivateKey, string(password))
				if err != nil {
					log.Fatalf("Error decrypting private key: %v", err)
				}
				wallet.PrivateKey = decryptedPK
			}

			// Decrypt Mnemonic if present
			if encWallet.Mnemonic != "" {
				decryptedMnemonic, err := encryption.Decrypt(encWallet.Mnemonic, string(password))
				if err != nil {
					log.Fatalf("Error decrypting mnemonic: %v", err)
				}
				wallet.Mnemonic = decryptedMnemonic
			}

			decryptedWallets = append(decryptedWallets, wallet)
		}

		// Create output directory if it doesn't exist
		outputDir := "output"
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			log.Fatalf("Error creating output directory: %v", err)
		}

		// Convert wallets to JSON (fully decrypted)
		jsonData, err := json.MarshalIndent(decryptedWallets, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling wallets: %v", err)
		}

		// Write decrypted JSON to output directory
		baseName := strings.TrimSuffix(filepath.Base(*decryptFile), ".encrypted.json")
		if baseName == filepath.Base(*decryptFile) {
			// If it doesn't have .encrypted.json suffix, just add .json
			baseName = strings.TrimSuffix(baseName, ".json")
		}
		outputFile := filepath.Join(outputDir, baseName+".json")
		if err := os.WriteFile(outputFile, jsonData, 0600); err != nil {
			log.Fatalf("Error writing decrypted file: %v", err)
		}

		fmt.Printf("Successfully decrypted %d wallets to %s\n", len(decryptedWallets), outputFile)
		return
	}

	// Wallet Address Validator
	r, err := regexp.Compile(*regEx)
	if err != nil {
		panic(err)
	}
	contains := strings.Split(*contain, ",")
	*prefix = utils.Add0xPrefix(*prefix)
	if strings.EqualFold(*engineName, "gpu") {
		args := strings.Fields(strings.TrimSpace(*gpuArgs))
		hasPrefixArg := hasAnyArg(args, "-p", "--prefix")
		hasSuffixArg := hasAnyArg(args, "-s", "--suffix")
		hasCreateArg := hasAnyArg(args, "-c", "--contract")
		hasCreate2Arg := hasAnyArg(args, "-c2", "--contract2")
		hasCreate3Arg := hasAnyArg(args, "-c3", "--contract3")
		hasBytecodeArg := hasAnyArg(args, "-b", "--bytecode")
		hasAddressArg := hasAnyArg(args, "-a", "--address")
		hasDeployerArg := hasAnyArg(args, "-da", "-ad", "--deployer-address")

		if *prefix != "" && !hasPrefixArg {
			args = append(args, "-p", strings.TrimPrefix(strings.TrimPrefix(*prefix, "0x"), "0X"))
		}
		if *suffix != "" && !hasSuffixArg {
			args = append(args, "-s", strings.TrimPrefix(strings.TrimPrefix(*suffix, "0x"), "0X"))
		}

		switch strings.ToLower(strings.TrimSpace(*contractMode)) {
		case "":
			// no-op
		case "create":
			if !hasCreateArg && !hasCreate2Arg && !hasCreate3Arg {
				args = append(args, "--contract")
			}
		case "create2":
			if !hasCreateArg && !hasCreate2Arg && !hasCreate3Arg {
				args = append(args, "--contract2")
			}
			if strings.TrimSpace(*contractBytecode) == "" {
				log.Fatal("missing -contract-bytecode for -contract-mode create2")
			}
			if strings.TrimSpace(*contractAddress) == "" {
				log.Fatal("missing -contract-address for -contract-mode create2")
			}
			if err := validateWorkerBytecodeFile(*contractBytecode); err != nil {
				log.Fatalf("invalid -contract-bytecode: %v", err)
			}
			if !hasBytecodeArg {
				args = append(args, "--bytecode", *contractBytecode)
			}
			if !hasAddressArg {
				args = append(args, "--address", *contractAddress)
			}
			// NOTE: upstream l3wi worker currently may also require --deployer-address
			// for --contract2 in some builds. Pass it when provided to avoid a hard stop.
			if strings.TrimSpace(*contractDeployer) != "" && !hasDeployerArg {
				args = append(args, "--deployer-address", *contractDeployer)
			}
		case "create3":
			if !hasCreateArg && !hasCreate2Arg && !hasCreate3Arg {
				args = append(args, "--contract3")
			}
			if strings.TrimSpace(*contractBytecode) == "" {
				log.Fatal("missing -contract-bytecode for -contract-mode create3")
			}
			if strings.TrimSpace(*contractAddress) == "" {
				log.Fatal("missing -contract-address for -contract-mode create3")
			}
			if strings.TrimSpace(*contractDeployer) == "" {
				log.Fatal("missing -contract-deployer for -contract-mode create3")
			}
			if err := validateWorkerBytecodeFile(*contractBytecode); err != nil {
				log.Fatalf("invalid -contract-bytecode: %v", err)
			}
			if !hasBytecodeArg {
				args = append(args, "--bytecode", *contractBytecode)
			}
			if !hasAddressArg {
				args = append(args, "--address", *contractAddress)
			}
			if !hasDeployerArg {
				args = append(args, "--deployer-address", *contractDeployer)
			}
		default:
			log.Fatalf("unsupported -contract-mode %q (use create|create2|create3)", *contractMode)
		}
		*gpuArgs = strings.Join(args, " ")
	}
	validateAddress := func(address string) bool {
		isValid := true
		if len(contains) > 0 {
			cb := func(contain string) bool {
				return strings.Contains(address, contain)
			}
			if *strict {
				if !utils.Have(contains, cb) {
					isValid = false
				}
			} else {
				if !utils.Some(contains, cb) {
					isValid = false
				}
			}
		}

		if *prefix != "" {
			if !strings.HasPrefix(address, *prefix) {
				isValid = false
			}
		}

		if *suffix != "" {
			if !strings.HasSuffix(address, *suffix) {
				isValid = false
			}
		}

		if *regEx != "" && !r.MatchString(address) {
			isValid = false
		}

		return isValid
	}
	if *number <= 0 {
		*number = -1
	}
	if *limit <= 0 {
		*limit = *number
	}

	// Progress bar
	var bar progressbar.ProgressBar
	if *isCompatible {
		bar = progressbar.NewCompatibleProgressBar(*number)
	} else {
		bar = progressbar.NewStandardProgressBar(*number)
	}

	// Repository
	var repo repository.Repository
	switch {
	case *dbPath != "":
		db, err := gorm.Open(sqlite.Open("./db/"+*dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Silent),
			DryRun:                 *isDryrun,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic(err)
		}

		defer func() {
			db, _ := db.DB()
			db.Close()
		}()

		if !*isDryrun {
			if err := db.AutoMigrate(&wallets.Wallet{}); err != nil {
				panic(err)
			}
		}

		repo = repository.NewGormRepository(db, uint64(*concurrency))
	default:
		repo = repository.NewInMemoryRepositoryWithConfig(repository.InMemoryConfig{
			OutputDir:  *outputDir,
			OutputName: *outputName,
			NoExport:   *noExport,
		})
	}

	engineInstance, err := enginefactory.New(enginefactory.Config{
		Mode:    *mode,
		Bits:    *bits,
		Engine:  *engineName,
		GPUBin:  *gpuBin,
		GPUArgs: *gpuArgs,
	})
	if err != nil {
		panic(err)
	}

	statusProvider := makeETAStatusProvider(*engineName, *prefix, *suffix, engineInstance.Speed)

	generator := generators.New(
		engineInstance.Generate,
		repo,
		generators.Config{
			AddresValidator: validateAddress,
			ProgressBar:     bar,
			SpeedProvider:   engineInstance.Speed,
			StatusProvider:  statusProvider,
			OutputMode:      *outputMode,
			DryRun:          *isDryrun,
			Concurrency:     *concurrency,
			Number:          *number,
			Limit:           *limit,
		},
	)

	go func() {
		<-ctx.Done()
		// Kill GPU worker first so any blocked Next() unblocks immediately.
		if err := engineInstance.Close(); err != nil {
			log.Printf("Engine Close Error: %+v", err)
		}

		if err := generator.Shutdown(); err != nil {
			log.Printf("Generator Shutdown Error: %+v", err)
		}
	}()

	if err := generator.Start(); err != nil {
		log.Printf("Generator Error: %+v", err)
	}
	if err := engineInstance.Close(); err != nil {
		log.Printf("Engine Close Error: %+v", err)
	}
	if err := repo.Close(); err != nil {
		log.Printf("WalletsRepo Close Error: %+v", err)
	}
}

func hasAnyArg(args []string, keys ...string) bool {
	for _, a := range args {
		for _, k := range keys {
			if a == k {
				return true
			}
		}
	}
	return false
}

func makeETAStatusProvider(engineName, prefix, suffix string, speedProvider func() string) func(resolved int64, elapsed time.Duration) string {
	difficultyNibbles := vanityDifficultyNibbles(prefix, suffix)
	return func(resolved int64, elapsed time.Duration) string {
		speedText := "n/a"
		speedWps := 0.0
		if strings.EqualFold(engineName, "gpu") && speedProvider != nil {
			if s := strings.TrimSpace(speedProvider()); s != "" {
				speedText = s
				speedWps = parseWorkerSpeedToWPS(s)
			}
		}

		eta := "n/a"
		if difficultyNibbles > 0 && speedWps > 0 {
			meanSec := math.Pow(16, float64(difficultyNibbles)) / speedWps
			p50 := math.Log(2) * meanSec
			p90 := -math.Log(0.1) * meanSec
			eta = fmt.Sprintf("exp=%s p50=%s p90=%s", formatDurationShort(time.Duration(meanSec*float64(time.Second))), formatDurationShort(time.Duration(p50*float64(time.Second))), formatDurationShort(time.Duration(p90*float64(time.Second))))
		}

		return fmt.Sprintf("resolved: %d | gpu: %s | elapsed: %s | eta(exp/p50/p90): %s", resolved, speedText, formatDurationShort(elapsed), eta)
	}
}

func vanityDifficultyNibbles(prefix, suffix string) int {
	p := strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(prefix), "0x"), "0X")
	s := strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(suffix), "0x"), "0X")
	if p == "" && s == "" {
		return 0
	}
	// Overlap case is ambiguous for exact probability in this simple model.
	if len(p)+len(s) > 40 {
		return 0
	}
	return len(p) + len(s)
}

func parseWorkerSpeedToWPS(speed string) float64 {
	// Expected input from worker parser: "6032.26 M/s"
	s := strings.TrimSpace(strings.ReplaceAll(speed, " ", ""))
	if !strings.HasSuffix(s, "M/s") {
		return 0
	}
	v := strings.TrimSuffix(s, "M/s")
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f <= 0 {
		return 0
	}
	return f * 1_000_000
}

func formatDurationShort(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// validateWorkerBytecodeFile validates bytecode for l3wi worker parser:
// even number of hex chars and no trailing whitespace/newline noise.
func validateWorkerBytecodeFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := strings.TrimSpace(string(data))
	content = strings.TrimPrefix(strings.TrimPrefix(content, "0x"), "0X")
	if content == "" {
		return fmt.Errorf("empty bytecode")
	}
	if len(content)%2 != 0 {
		return fmt.Errorf("hex length must be even")
	}
	for _, r := range content {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return fmt.Errorf("contains non-hex character %q", r)
		}
	}
	return nil
}
