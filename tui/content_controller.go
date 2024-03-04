package tui

import (
	"context"
	"strings"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/pack/command"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type controllerEvent int16

const (
	eventResetFilter controllerEvent = iota
	eventActivateFilter
	eventCancelFilter
	eventFinishFilter
	eventRefreshContentByFilter
	eventShowDescription
	eventHideDescription
)

type contentController struct {
	events chan controllerEvent
}

func newContentController(ctx context.Context, views *appViews, commandEntities []command.CommandEntity) *contentController {
	initialCommands := mapCommandsToViewItems(commandEntities)
	events := make(chan controllerEvent, 100)
	contentState := contentController{events: events}
	go processControlsEvents(ctx, views, initialCommands, events)
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

func (cs *contentController) ShowDescription() {
	cs.events <- eventShowDescription
}

func (cs *contentController) HideDescription() {
	cs.events <- eventHideDescription
}

func processControlsEvents(ctx context.Context, views *appViews, initialCommands []commandViewItem, events chan controllerEvent) {
	app := application.NewApplication(application.NewController(views.app))
	populateCommandsView(ctx, views.commands, app, initialCommands)

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
				if views.commands.GetCurrentItem() != 0 {
					views.commands.SetCurrentItem(0)
					views.app.Draw()
				}
				return
			}
			currentText = ""
			views.filter.SetText("", true)
			absoluteCommandIndex := getAbsoluteCommandIndex(views.commands, filteredCommands)
			filteredCommands = initialCommands
			views.commands.Clear()
			populateCommandsView(ctx, views.commands, app, filteredCommands)
			views.commands.SetCurrentItem(absoluteCommandIndex)
			views.app.Draw()
		case eventActivateFilter:
			if isFilterViewActive {
				return
			}
			isFilterViewActive = true
			views.filter.SetDisabled(false)
			views.app.SetFocus(views.filter)
			views.app.Draw()
		case eventCancelFilter:
			if !isFilterViewActive {
				return
			}
			isFilterViewActive = false
			currentText = ""
			views.filter.SetDisabled(true)
			views.filter.SetText("", true)
			absoluteCommandIndex := getAbsoluteCommandIndex(views.commands, filteredCommands)
			filteredCommands = initialCommands
			views.commands.Clear()
			populateCommandsView(ctx, views.commands, app, filteredCommands)
			views.commands.SetCurrentItem(absoluteCommandIndex)
			views.app.SetFocus(views.commands)
			views.app.Draw()
		case eventFinishFilter:
			if !isFilterViewActive {
				return
			}
			isFilterViewActive = false
			views.filter.SetDisabled(true)
			views.app.SetFocus(views.commands)
			views.app.Draw()
		case eventRefreshContentByFilter:
			newText := views.filter.GetText()
			if currentText == newText {
				return
			}
			currentText = newText
			absoluteCommandIndex := getAbsoluteCommandIndex(views.commands, filteredCommands)
			commands, index := filterCommands(initialCommands, absoluteCommandIndex, newText)
			filteredCommands = commands
			views.commands.Clear()
			populateCommandsView(ctx, views.commands, app, commands)
			views.commands.SetCurrentItem(index)
		case eventShowDescription:
			item := getCurrentCommandViewItem(views.commands, filteredCommands)
			if item == nil {
				return
			}
			views.pages.AddAndSwitchToPage(pageNameDescription, createDescriptionPage(item), true)
			views.app.Draw()
		case eventHideDescription:
			views.pages.RemovePage(pageNameDescription)
			views.app.Draw()
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

func getCurrentCommandViewItem(commandsView *tview.List, commands []commandViewItem) *commandViewItem {
	currentItem := commandsView.GetCurrentItem()
	if currentItem >= 0 && len(commands) > 0 {
		return &commands[currentItem]
	}
	return nil
}

func getAbsoluteCommandIndex(commandsView *tview.List, commands []commandViewItem) int {
	if item := getCurrentCommandViewItem(commandsView, commands); item != nil {
		return item.Index
	}
	return 0
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

func createDescriptionPage(item *commandViewItem) tview.Primitive {
	descriptionView := tview.NewTextView().SetText(item.CommandEntity.Properties.Name)
	frame := tview.NewFrame(descriptionView)
	frame.SetBorders(1, 0, 1, 1, 1, 1).
		SetBorder(true).
		SetTitle("Description")
	frame.AddText(item.CommandEntity.Properties.Name, true, tview.AlignCenter, tcell.ColorWhite)
	if item.CommandEntity.Properties.Alias != "" {
		frame.AddText("alias: "+item.CommandEntity.Properties.Alias, true, tview.AlignCenter, tcell.ColorWhite)
	}
	frame.AddText("esc", false, tview.AlignRight, tcell.ColorWhite)
	return frame
}
