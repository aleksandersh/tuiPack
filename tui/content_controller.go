package tui

import (
	"context"
	"strings"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/pack/command"
	"github.com/rivo/tview"
)

type controllerEvent int16

const (
	eventResetFilter controllerEvent = iota
	eventActivateFilter
	eventCancelFilter
	eventFinishFilter
	eventRefreshContentByFilter
)

type contentController struct {
	events chan controllerEvent
}

func newContentController(ctx context.Context, app *tview.Application, commandsView *tview.List, filterView *tview.TextArea, commandEntities []command.CommandEntity) *contentController {
	initialCommands := mapCommandsToViewItems(commandEntities)
	events := make(chan controllerEvent, 100)
	contentState := contentController{events: events}
	go processControlsEvents(ctx, app, commandsView, filterView, initialCommands, events)
	return &contentState
}

func (cs *contentController) ResetFilter() {
	cs.events <- eventResetFilter
}

func (cs *contentController) ActivateFilter() {
	cs.events <- eventActivateFilter
}

func (cs *contentController) CancelFilter() {
	cs.events <- eventCancelFilter
}

func (cs *contentController) FinishFilter() {
	cs.events <- eventFinishFilter
}

func (cs *contentController) RefreshContentByFilter() {
	cs.events <- eventRefreshContentByFilter
}

func processControlsEvents(ctx context.Context, tviewApp *tview.Application, commandsView *tview.List, filterView *tview.TextArea, initialCommands []commandViewItem, events chan controllerEvent) {
	app := application.NewApplication(application.NewController(tviewApp))
	populateCommandsView(ctx, commandsView, app, initialCommands)

	filteredCommands := initialCommands
	isFilterViewActive := false
	currentText := ""

	observeEvents(ctx, events, func(event controllerEvent) {
		switch event {
		case eventResetFilter:
			if isFilterViewActive {
				return
			}
			if currentText == "" {
				if commandsView.GetCurrentItem() != 0 {
					commandsView.SetCurrentItem(0)
					tviewApp.Draw()
				}
				return
			}
			currentText = ""
			filterView.SetText("", true)
			absoluteCommandIndex := getAbsoluteCommandIndex(commandsView, filteredCommands)
			filteredCommands = initialCommands
			commandsView.Clear()
			populateCommandsView(ctx, commandsView, app, filteredCommands)
			commandsView.SetCurrentItem(absoluteCommandIndex)
			tviewApp.Draw()
		case eventActivateFilter:
			if isFilterViewActive {
				return
			}
			isFilterViewActive = true
			filterView.SetDisabled(false)
			tviewApp.SetFocus(filterView)
			tviewApp.Draw()
		case eventCancelFilter:
			if !isFilterViewActive {
				return
			}
			isFilterViewActive = false
			currentText = ""
			filterView.SetDisabled(true)
			filterView.SetText("", true)
			absoluteCommandIndex := getAbsoluteCommandIndex(commandsView, filteredCommands)
			filteredCommands = initialCommands
			commandsView.Clear()
			populateCommandsView(ctx, commandsView, app, filteredCommands)
			commandsView.SetCurrentItem(absoluteCommandIndex)
			tviewApp.SetFocus(commandsView)
			tviewApp.Draw()
		case eventFinishFilter:
			if !isFilterViewActive {
				return
			}
			isFilterViewActive = false
			filterView.SetDisabled(true)
			tviewApp.SetFocus(commandsView)
			tviewApp.Draw()
		case eventRefreshContentByFilter:
			newText := filterView.GetText()
			if currentText == newText {
				return
			}
			currentText = newText
			absoluteCommandIndex := getAbsoluteCommandIndex(commandsView, filteredCommands)
			commands, index := filterCommands(initialCommands, absoluteCommandIndex, newText)
			filteredCommands = commands
			commandsView.Clear()
			populateCommandsView(ctx, commandsView, app, commands)
			commandsView.SetCurrentItem(index)
		}
	})
}

func observeEvents(ctx context.Context, events chan controllerEvent, handler func(event controllerEvent)) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-events:
			handler(event)
		}
	}
}

func getAbsoluteCommandIndex(commandsView *tview.List, commands []commandViewItem) int {
	currentItem := commandsView.GetCurrentItem()
	absoluteFocusedIndex := 0
	if currentItem >= 0 && len(commands) > 0 {
		absoluteFocusedIndex = commands[currentItem].Index
	}
	return absoluteFocusedIndex
}

func filterCommands(items []commandViewItem, absoluteFocusedIndex int, text string) ([]commandViewItem, int) {
	filteredItems := make([]commandViewItem, 0, len(items))
	focusedIndex := 0
	for _, item := range items {
		if strings.Contains(item.Text, text) {
			if item.Index == absoluteFocusedIndex {
				focusedIndex = len(filteredItems)
			}
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems, focusedIndex
}

func populateCommandsView(ctx context.Context, commandsView *tview.List, app *application.Application, items []commandViewItem) {
	for _, item := range items {
		addCommandView(ctx, commandsView, app, item)
	}
}

func addCommandView(ctx context.Context, listView *tview.List, app *application.Application, item commandViewItem) {
	properties := item.CommandEntity.Properties
	listView.AddItem(properties.Name, properties.Description, 0, func() {
		item.CommandEntity.Command.Execute(ctx, app, properties)
	})
}
