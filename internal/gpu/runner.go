package gpu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/planxnx/ethereum-wallet-generator/wallets"
)

var legacyGPULineRE = regexp.MustCompile(`Private Key:\s*0x([0-9a-fA-F]{64}).*Address:\s*0x([0-9a-fA-F]{40})`)

type Runner struct {
	cmd       *exec.Cmd
	out       io.ReadCloser
	walletCh  chan *wallets.Wallet
	errCh     chan error
	closeOnce sync.Once
	stderrMu  sync.Mutex
	stderrLog []string
}

// NewRunner starts an external GPU worker process that continuously prints wallets to stdout.
// Supported output formats (one wallet per line):
// 1) JSON: {"address":"0x...","privateKey":"..."} (private_key also accepted)
// 2) Plain text: "<address> <privateKey>"
func NewRunner(bin string, args []string) (*Runner, error) {
	cmd := exec.Command(bin, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("gpu stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("gpu stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start gpu process: %w", err)
	}

	r := &Runner{
		cmd:      cmd,
		out:      stdout,
		walletCh: make(chan *wallets.Wallet, 256),
		errCh:    make(chan error, 2),
	}

	go r.readStdout()
	go r.readStderr(stderr)
	go r.wait()

	return r, nil
}

func (r *Runner) Next() (*wallets.Wallet, error) {
	select {
	case w, ok := <-r.walletCh:
		if !ok {
			// Give wait() a short window to publish the real process exit error.
			select {
			case err := <-r.errCh:
				if err != nil {
					return nil, fmt.Errorf("gpu wallet stream closed: %w%s", err, r.stderrSuffix())
				}
			case <-time.After(250 * time.Millisecond):
			}
			return nil, fmt.Errorf("gpu wallet stream closed%s", r.stderrSuffix())
		}
		return w, nil
	case err := <-r.errCh:
		if err == nil {
			return nil, fmt.Errorf("gpu process exited%s", r.stderrSuffix())
		}
		return nil, fmt.Errorf("%w%s", err, r.stderrSuffix())
	}
}

func (r *Runner) Close() error {
	r.closeOnce.Do(func() {
		if r.cmd.Process != nil {
			_ = r.cmd.Process.Kill()
		}
	})
	return nil
}

func (r *Runner) readStdout() {
	sc := bufio.NewScanner(r.out)
	sc.Split(scanCRLF)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		w, err := parseWalletLine(line)
		if err != nil {
			continue
		}

		r.walletCh <- w
	}
	if err := sc.Err(); err != nil {
		r.pushErr(fmt.Errorf("gpu stdout scan: %w", err))
	}
	close(r.walletCh)
}

// scanCRLF splits on either '\n' or '\r'. This is needed because some GPU
// workers print progress with carriage returns (no newline).
func scanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		// Return line up to separator and consume one separator byte.
		return i + 1, data[:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func (r *Runner) readStderr(stderr io.ReadCloser) {
	// Drain stderr to prevent child process from blocking when stderr buffer is full.
	sc := bufio.NewScanner(stderr)
	for sc.Scan() {
		r.stderrMu.Lock()
		r.stderrLog = append(r.stderrLog, strings.TrimSpace(sc.Text()))
		if len(r.stderrLog) > 8 {
			r.stderrLog = r.stderrLog[len(r.stderrLog)-8:]
		}
		r.stderrMu.Unlock()
	}
	if err := sc.Err(); err != nil {
		r.pushErr(fmt.Errorf("gpu stderr scan: %w", err))
	}
}

func (r *Runner) stderrSuffix() string {
	r.stderrMu.Lock()
	defer r.stderrMu.Unlock()
	if len(r.stderrLog) == 0 {
		return ""
	}
	return fmt.Sprintf(" (stderr: %s)", strings.Join(r.stderrLog, " | "))
}

func (r *Runner) wait() {
	r.pushErr(r.cmd.Wait())
}

func (r *Runner) pushErr(err error) {
	select {
	case r.errCh <- err:
	default:
	}
}

func parseWalletLine(line string) (*wallets.Wallet, error) {
	// Support legacy worker output, e.g.
	// "Elapsed: 000001 Score: 12 Private Key: 0x... Address: 0x..."
	if m := legacyGPULineRE.FindStringSubmatch(line); len(m) == 3 {
		address := normalizeAddress(m[2])
		privateKey := normalizePrivateKey(m[1])
		if common.IsHexAddress(address) && privateKey != "" {
			return &wallets.Wallet{
				Address:    address,
				PrivateKey: privateKey,
			}, nil
		}
	}

	// Try JSON first.
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err == nil {
		address := normalizeAddress(strFromAny(m["address"]))
		privateKey := normalizePrivateKey(strFromAny(m["privateKey"]))
		if privateKey == "" {
			privateKey = normalizePrivateKey(strFromAny(m["private_key"]))
		}

		if !common.IsHexAddress(address) || privateKey == "" {
			return nil, fmt.Errorf("invalid wallet json line")
		}

		return &wallets.Wallet{
			Address:    address,
			PrivateKey: privateKey,
		}, nil
	}

	// Fallback plain text: "<address> <privateKey>"
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid gpu line")
	}

	address := normalizeAddress(parts[0])
	privateKey := normalizePrivateKey(parts[1])
	if !common.IsHexAddress(address) || privateKey == "" {
		return nil, fmt.Errorf("invalid plain wallet line")
	}

	return &wallets.Wallet{
		Address:    address,
		PrivateKey: privateKey,
	}, nil
}

func strFromAny(v any) string {
	s, _ := v.(string)
	return s
}

func normalizeAddress(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if !strings.HasPrefix(strings.ToLower(s), "0x") {
		return "0x" + s
	}
	return s
}

func normalizePrivateKey(s string) string {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X"))
	if len(s) != 64 {
		return ""
	}
	return s
}
