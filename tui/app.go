package tui

import (
	"fmt"

	"github.com/aleksandersh/tuiPack/config"
	"github.com/aleksandersh/tuiPack/executor"
	"github.com/mattn/go-shellwords"
	"github.com/rivo/tview"
)

func RunApp(config *config.Pack) error {
	app := tview.NewApplication()
	contentView := createContentView(app, config)

	app.SetRoot(contentView, true).SetFocus(contentView)
	if err := app.Run(); err != nil {
		return fmt.Errorf("error in app.Run: %w", err)
	}
	app.Stop()
	return nil
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
		executor.ExecuteCommand(command.Args)
	})
}
