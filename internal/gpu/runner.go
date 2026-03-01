package gpu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/planxnx/ethereum-wallet-generator/wallets"
)

type Runner struct {
	cmd       *exec.Cmd
	out       io.ReadCloser
	walletCh  chan *wallets.Wallet
	errCh     chan error
	closeOnce sync.Once
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
			return nil, fmt.Errorf("gpu wallet stream closed")
		}
		return w, nil
	case err := <-r.errCh:
		if err == nil {
			return nil, fmt.Errorf("gpu process exited")
		}
		return nil, err
	}
}

func (r *Runner) Close() error {
	var err error
	r.closeOnce.Do(func() {
		if r.cmd.Process != nil {
			_ = r.cmd.Process.Kill()
		}
		err = r.cmd.Wait()
	})
	return err
}

func (r *Runner) readStdout() {
	sc := bufio.NewScanner(r.out)
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
	close(r.walletCh)
}

func (r *Runner) readStderr(stderr io.ReadCloser) {
	// Drain stderr to prevent child process from blocking when stderr buffer is full.
	sc := bufio.NewScanner(stderr)
	for sc.Scan() {
	}
}

func (r *Runner) wait() {
	r.errCh <- r.cmd.Wait()
}

func parseWalletLine(line string) (*wallets.Wallet, error) {
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

