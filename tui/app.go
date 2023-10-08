package tui

import (
	"context"
	"fmt"

	"github.com/aleksandersh/tuiPack/config"
	"github.com/aleksandersh/tuiPack/executor"
	"github.com/rivo/tview"
)

func RunApp(ctx context.Context, config *config.Pack) error {
	app := tview.NewApplication()
	contentView := createContentView(ctx, app, config)

	app.SetRoot(contentView, true).SetFocus(contentView)
	if err := app.Run(); err != nil {
		return fmt.Errorf("error in app.Run: %w", err)
	}
	app.Stop()
	return nil
}

func createContentView(ctx context.Context, app *tview.Application, config *config.Pack) tview.Primitive {
	commandsView := tview.NewList()
	commandsView.SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetWrapAround(false).
		SetTitle(config.Name).
		SetBorder(true)

	for _, command := range config.Commands {
		addCommandView(ctx, app, commandsView, command)
	}

	commandsView.Focus(func(p tview.Primitive) {})

	return commandsView
}

func addCommandView(ctx context.Context, app *tview.Application, listView *tview.List, command config.Command) {
	listView.AddItem(command.Name, command.Description, 0, func() {
		app.Stop()
		executor.ExecuteCommand(ctx, command.Args, command.Environment)
	})
}
