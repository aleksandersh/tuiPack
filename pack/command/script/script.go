package script

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/pack/command"
	"github.com/mattn/go-shellwords"
)

type scriptCommand struct {
	script string
}

func newScript(script string) *scriptCommand {
	return &scriptCommand{script: script}
}

func (cmd *scriptCommand) Execute(ctx context.Context, app *application.Application, props *command.Properties) {
	app.Ui.Close()

	args, err := cmd.parseScript(props)
	if err != nil {
		log.Fatalf("error in script.parseScript: %v", err)
		return
	}

	env := cmd.prepareEnv(props)

	execute(ctx, args, env)
}

func (cmd *scriptCommand) parseScript(props *command.Properties) ([]string, error) {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	parser.Getenv = func(env string) string {
		if value, contains := props.Environment[env]; contains {
			return value
		}
		return os.Getenv(env)
	}
	args, err := parser.Parse(cmd.script)
	if err != nil {
		return []string{}, fmt.Errorf("error in parser.Parse: %w", err)
	}
	return args, nil
}

func (cmd *scriptCommand) prepareEnv(props *command.Properties) []string {
	env := make([]string, 0, len(props.Environment))
	for key, value := range props.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
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
