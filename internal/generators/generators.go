package generators

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/planxnx/ethereum-wallet-generator/internal/progressbar"
	"github.com/planxnx/ethereum-wallet-generator/internal/repository"
	"github.com/planxnx/ethereum-wallet-generator/wallets"
)

type Config struct {
	AddresValidator func(address string) bool
	ProgressBar     progressbar.ProgressBar
	SpeedProvider   func() string
	StatusProvider  func(resolved int64, elapsed time.Duration) string
	OutputMode      string
	DryRun          bool
	Concurrency     int
	Number          int
	Limit           int
}

type Generator struct {
	walletGen wallets.Generator
	repo      repository.Repository
	config    Config

	isShutdown     atomic.Bool
	shutdownSignal chan struct{}
	shutdownOnce   sync.Once
	shutDownWg     sync.WaitGroup
}

func New(walletGen wallets.Generator, repo repository.Repository, cfg Config) *Generator {
	return &Generator{
		walletGen:      walletGen,
		repo:           repo,
		config:         cfg,
		shutdownSignal: make(chan struct{}),
	}
}

func (g *Generator) requestShutdown() {
	g.shutdownOnce.Do(func() {
		close(g.shutdownSignal)
	})
}

func (g *Generator) Start() (err error) {
	g.isShutdown.Store(false)
	g.shutDownWg.Add(1)
	defer g.shutDownWg.Done()

	var (
		bar           = g.config.ProgressBar
		resolvedCount atomic.Int64
		start         = time.Now()
	)
	defer func() {
		_ = bar.Finish()

		if w := g.repo.Result(); len(w) > 0 && !g.config.DryRun {
			switch strings.ToLower(strings.TrimSpace(g.config.OutputMode)) {
			case "full":
				col2Name := "Seed"
				var result strings.Builder
				for _, wallet := range w {
					col2 := wallet.Mnemonic
					if wallet.Mnemonic == "" {
						col2 = wallet.PrivateKey
						col2Name = "Private Key"
					}
					if _, err := fmt.Fprintf(&result, "%-42s %s\n", wallet.Address, col2); err != nil {
						continue
					}
				}
				fmt.Printf("\n%-42s %s\n", "Address", col2Name)
				fmt.Printf("%-42s %s\n", strings.Repeat("-", 42), strings.Repeat("-", 90))
				fmt.Println(result.String())
			case "silent":
				// no terminal wallet output
			default:
				// safe: no secret material in terminal
				var result strings.Builder
				for _, wallet := range w {
					if _, err := fmt.Fprintf(&result, "%-42s %s\n", wallet.Address, "<hidden>"); err != nil {
						continue
					}
				}
				fmt.Printf("\n%-42s %s\n", "Address", "Secret")
				fmt.Printf("%-42s %s\n", strings.Repeat("-", 42), strings.Repeat("-", 90))
				fmt.Println(result.String())
			}
		}

		fmt.Printf("\nResolved Speed: %.2f w/s\n", float64(resolvedCount.Load())/time.Since(start).Seconds())
		fmt.Printf("Total Duration: %v\n", time.Since(start))
		fmt.Printf("Total Wallet Resolved: %d w\n", resolvedCount.Load())
		fmt.Printf("\nCopyright (C) 2023 Planxnx <planxthanee@gmail.com>\n")

		g.isShutdown.Store(true)
	}()

	var wg sync.WaitGroup
	var statusWG sync.WaitGroup
	commands := make(chan struct{})
	updateStatus := func() {
		if g.config.StatusProvider != nil {
			_ = bar.SetStatus(g.config.StatusProvider(resolvedCount.Load(), time.Since(start)))
			return
		}
		speed := "n/a"
		if g.config.SpeedProvider != nil {
			if s := strings.TrimSpace(g.config.SpeedProvider()); s != "" {
				speed = s
			}
		}
		_ = bar.SetStatus(fmt.Sprintf("resolved: %d | gpu: %s | elapsed: %s | eta(exp/p50/p90): n/a", resolvedCount.Load(), speed, time.Since(start).Round(time.Second)))
	}
	statusWG.Add(1)
	go func() {
		defer statusWG.Done()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-g.shutdownSignal:
				return
			case <-ticker.C:
				updateStatus()
			}
		}
	}()
	for i := 0; i < g.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range commands {
				if !(resolvedCount.Load() < int64(g.config.Limit) || g.config.Limit < 0) {
					return
				}

				wallet, err := g.walletGen()
				if err != nil {
					// Stop early on fatal GPU worker termination to avoid log floods.
					if strings.Contains(err.Error(), "gpu wallet stream closed") || strings.Contains(err.Error(), "gpu process exited") {
						log.Printf("[ERROR] fatal generator error: %+v\n", err)
						g.requestShutdown()
						return
					}
					// Ignore non-fatal error
					log.Printf("[ERROR] failed to generate wallet: %+v\n", err)
					continue
				}

				isOk := true
				if g.config.AddresValidator != nil {
					isOk = g.config.AddresValidator(wallet.Address)
				}

				if isOk {
					if err := g.repo.Insert(wallet); err != nil {
						// Ignore error
						log.Printf("[ERROR] failed to insert wallet to db: %+v\n", err)
						continue
					}
					resolvedCount.Add(1)
				}

				_ = bar.SetResolved(int(resolvedCount.Load()))
				updateStatus()
				_ = bar.Increment()
			}
		}()
	}

mainloop:
	for i := 0; (i < g.config.Number || g.config.Number < 0) && (resolvedCount.Load() < int64(g.config.Limit) || g.config.Limit < 0); i++ {
		select {
		case <-g.shutdownSignal:
			break mainloop
		case commands <- struct{}{}:
			// submitted
		}
	}

	close(commands)
	wg.Wait()
	g.requestShutdown()
	statusWG.Wait()
	return nil
}

func (g *Generator) Shutdown() (err error) {
	if g.isShutdown.Load() {
		return nil
	}
	g.requestShutdown()
	g.shutDownWg.Wait()
	return nil
}
