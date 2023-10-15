package tui

import (
	"context"
	"fmt"

	"github.com/aleksandersh/tuiPack/app/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	keySlash = 47
)

func RunApp(ctx context.Context, config *config.Pack) error {
	app := tview.NewApplication()

	commandsView := createCommandsView(ctx, app, config)
	filterView := createFilterView()
	containerView := createContainerView(commandsView, filterView)

	setupContent(ctx, app, commandsView, filterView, config.Commands)

	app.SetRoot(containerView, true).SetFocus(commandsView)
	if err := app.Run(); err != nil {
		return fmt.Errorf("error in app.Run: %w", err)
	}
	app.Stop()
	return nil
}

func createCommandsView(ctx context.Context, app *tview.Application, config *config.Pack) *tview.List {
	commandsView := tview.NewList()
	commandsView.SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetWrapAround(false).
		SetTitle(config.Name).
		SetBorder(true)

	return commandsView
}

func createFilterView() *tview.TextArea {
	filterView := tview.NewTextArea()
	filterView.SetDisabled(true)
	filterView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune || event.Key() == tcell.KeyBackspace2 {
			return event
		}
		return nil
	})
	return filterView
}

func createContainerView(commandsView tview.Primitive, filterView tview.Primitive) *tview.Flex {
	containerView := tview.NewFlex()
	containerView.
		SetDirection(tview.FlexRow).
		AddItem(commandsView, 0, 1, true).
		AddItem(filterView, 1, 0, false)
	return containerView
}

func setupContent(ctx context.Context, app *tview.Application, commandsView *tview.List, filterView *tview.TextArea, commands []config.Command) {
	contentController := newContentController(ctx, app, commandsView, filterView, commands)
	isFilterViewActive := false
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !isFilterViewActive {
			if event.Key() == tcell.KeyRune && event.Rune() == keySlash {
				contentController.ActivateFilter()
				isFilterViewActive = true
				return nil
			}
			if event.Key() == tcell.KeyEsc {
				contentController.ResetFilter()
				return nil
			}
		}
		if !isFilterViewActive && event.Key() == tcell.KeyRune && event.Rune() == keySlash {
			contentController.ActivateFilter()
			isFilterViewActive = true
			return nil
		}
		if isFilterViewActive {
			if event.Key() == tcell.KeyEsc {
				contentController.CancelFilter()
				isFilterViewActive = false
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				contentController.FinishFilter()
				isFilterViewActive = false
				return nil
			}
		}
		return event
	})

	filterView.SetChangedFunc(func() {
		contentController.RefreshContentByFilter()
	})
}
