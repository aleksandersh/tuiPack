package executor

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func ExecuteCommand(ctx context.Context, args []string, env []string) {
	name := args[0]
	args = args[1:]
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(env) > 0 {
		cmd.Env = append(cmd.Environ(), env...)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	if err := cmd.Start(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Start: %v", err)
		}
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case signal := <-signals:
				cmd.Process.Signal(signal)
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Wait: %v", err)
		}
	}
}
