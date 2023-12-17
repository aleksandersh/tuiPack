package tui

import (
	"strings"

	"github.com/aleksandersh/tuiPack/command"
)

type commandViewItem struct {
	Index             int
	Text              string
	CommandProperties *command.Properties
	Command           command.Command
}

func mapCommandsToViewItems(commands []command.Command) []commandViewItem {
	items := make([]commandViewItem, 0, len(commands))
	for index, command := range commands {
		properties := command.GetProperties()
		items = append(items, commandViewItem{
			Index:             index,
			Text:              strings.ToLower(properties.Name),
			CommandProperties: properties,
			Command:           command,
		})
	}
	return items
}
