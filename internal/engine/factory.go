package engine

import (
	"fmt"
	"strings"

	"github.com/planxnx/ethereum-wallet-generator/internal/gpu"
	"github.com/planxnx/ethereum-wallet-generator/wallets"
)

type Config struct {
	Mode      int
	Bits      int
	Engine    string
	GPUBin    string
	GPUArgs   string
}

type Instance struct {
	Generate wallets.Generator
	closeFn  func() error
	speedFn  func() string
}

func (i *Instance) Close() error {
	if i == nil || i.closeFn == nil {
		return nil
	}
	return i.closeFn()
}

func (i *Instance) Speed() string {
	if i == nil || i.speedFn == nil {
		return ""
	}
	return i.speedFn()
}

func New(cfg Config) (*Instance, error) {
	engineName := strings.ToLower(strings.TrimSpace(cfg.Engine))
	if engineName == "" {
		engineName = "cpu"
	}

	switch cfg.Mode {
	case 1:
		if engineName != "cpu" {
			return nil, fmt.Errorf("gpu engine currently supports only -mode 2")
		}
		return &Instance{
			Generate: wallets.NewGeneratorMnemonic(cfg.Bits),
		}, nil
	case 2:
		if engineName == "gpu" {
			if strings.TrimSpace(cfg.GPUBin) == "" {
				// Default to bundled worker name available in GPU images.
				cfg.GPUBin = "vanity-eth-address"
			}
			args := strings.Fields(strings.TrimSpace(cfg.GPUArgs))
			runner, err := gpu.NewRunner(cfg.GPUBin, args)
			if err != nil {
				return nil, err
			}
			return &Instance{
				Generate: runner.Next,
				closeFn:  runner.Close,
				speedFn:  runner.LatestSpeed,
			}, nil
		}
		return &Instance{
			Generate: wallets.NewGeneratorPrivatekey(),
		}, nil
	default:
		return nil, fmt.Errorf("invalid mode. See: https://github.com/Planxnx/ethereum-wallet-generator#Modes")
	}
}
