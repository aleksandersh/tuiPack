package tui

import (
	"strings"

	"github.com/aleksandersh/tuiPack/pack/command"
)

type commandViewItem struct {
	Index         int
	Text          string
	CommandEntity command.CommandEntity
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
