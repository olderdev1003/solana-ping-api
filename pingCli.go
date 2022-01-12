package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func solanaPing(c Cluster) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var configpath string
	switch c {
	case MainnetBeta:
		configpath = config.SolanaConfigPath + config.SolanaConfig.Mainnet
	case Testnet:
		configpath = config.SolanaConfigPath + config.SolanaConfig.Testnet
	case Devnet:
		configpath = config.SolanaConfigPath + config.SolanaConfig.Devnet
	default:
		configpath = config.SolanaConfigPath + config.SolanaConfig.Devnet
	}
	log.Info("configfile=", configpath)
	cmd := exec.CommandContext(ctx, "solana", "ping",
		"-c", fmt.Sprintf("%d", config.SolanaPing.Count),
		"-i", fmt.Sprintf("%d", config.SolanaPing.Interval),
		"-C", configpath)
	cmd.Env = append(os.Environ(), ":"+config.SolanaPing.PingExePath)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
