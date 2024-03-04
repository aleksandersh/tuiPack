package tui

import (
	"strings"

	"github.com/aleksandersh/tuiPack/pack/command"
	"github.com/rivo/tview"
)

const (
	pageNameContent     = "content"
	pageNameDescription = "description"
)

type appViews struct {
	app      *tview.Application
	pages    *tview.Pages
	commands *tview.List
	filter   *tview.TextArea
}

type commandViewItem struct {
	Index         int
	Text          string
	CommandEntity command.CommandEntity
}

func newAppViews(app *tview.Application, pagesView *tview.Pages, commandsView *tview.List, filterView *tview.TextArea) *appViews {
	return &appViews{
		app:      app,
		pages:    pagesView,
		commands: commandsView,
		filter:   filterView,
	}
}

func mapCommandsToViewItems(commandEntities []command.CommandEntity) []commandViewItem {
	items := make([]commandViewItem, 0, len(commandEntities))
	for index, entity := range commandEntities {
		items = append(items, commandViewItem{
			Index:         index,
			Text:          strings.ToLower(entity.Properties.Name),
			CommandEntity: entity,
		})
	}
	return items
}
