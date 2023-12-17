package script

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/command"
)

type Script struct {
	properties *command.Properties
	args       []string
	env        []string
}

func newScript(properties *command.Properties, args []string, env []string) *Script {
	return &Script{properties: properties, args: args, env: env}
}

func (script Script) GetProperties() *command.Properties {
	return script.properties
}

func (script Script) Execute(ctx context.Context, app *application.Application) {
	app.Ui.Close()
	execute(ctx, script.args, script.env)
}

func execute(ctx context.Context, args []string, env []string) {
	name := args[0]
	args = args[1:]

	cmd := createScriptCommand(ctx, name, args, env)

	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	startScriptAsync(cmd)
	handleSystemSignalsAsync(ctx, signals, cmd)
	awaitForScriptCompletion(cmd)
}

func createScriptCommand(ctx context.Context, name string, args []string, env []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if len(env) > 0 {
		cmd.Env = append(cmd.Environ(), env...)
	}
	return cmd
}

func startScriptAsync(cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Start: %v", err)
		}
	}
}

func awaitForScriptCompletion(cmd *exec.Cmd) {
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Wait: %v", err)
		}
	}
}

func handleSystemSignalsAsync(ctx context.Context, signals chan os.Signal, cmd *exec.Cmd) {
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
}
