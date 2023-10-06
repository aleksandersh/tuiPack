package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/aleksandersh/tuiPack/config"
	"github.com/mattn/go-shellwords"
	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing config argument")
	}
	configFile := os.Args[1]
	config, err := config.ReadConfigFromYamlFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	app := tview.NewApplication()
	contentView := createContentView(app, config)

	app.SetRoot(contentView, true).SetFocus(contentView)
	if err := app.Run(); err != nil {
		log.Fatalf("failed to run tui application: %v", err)
	}
	app.Stop()
}

func createContentView(app *tview.Application, config *config.Pack) tview.Primitive {
	parser := shellwords.NewParser()
	parser.ParseEnv = true

	commandsView := tview.NewList()

	commandsView.SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetWrapAround(false).
		SetTitle(config.Name).
		SetBorder(true)

	for _, command := range config.Commands {
		addCommandView(app, parser, commandsView, command)
	}

	commandsView.Focus(func(p tview.Primitive) {})

	return commandsView
}

func addCommandView(app *tview.Application, parser *shellwords.Parser, listView *tview.List, command config.Command) {
	listView.AddItem(command.Name, command.Description, 0, func() {
		app.Stop()
		executeCommand(command.Args)
	})
}

func executeCommand(args []string) {
	name := args[0]
	args = args[1:]
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Run: %v", err)
		}
	}
}
