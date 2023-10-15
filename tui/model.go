package tui

import (
	"strings"

	"github.com/aleksandersh/tuiPack/app/config"
)

type commandViewItem struct {
	Index   int
	Text    string
	Command config.Command
}

func mapCommandsToViewItems(commands []config.Command) []commandViewItem {
	items := make([]commandViewItem, 0, len(commands))
	for index, command := range commands {
		items = append(items, commandViewItem{
			Index:   index,
			Text:    strings.ToLower(command.Name),
			Command: command,
		})
	}
	return items
}
