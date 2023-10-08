package command

import (
	"fmt"
	"strings"

	"github.com/aleksandersh/tuiPack/app"
	"github.com/aleksandersh/tuiPack/app/config"
)

const (
	leftColumnAliasLength = 30
)

func PrintAliases(args *app.Args, config *config.Pack) {
	fmt.Printf("Usage: tuiPack --script ALIAS\n\n")
	fmt.Printf("Aliases for available scripts:\n")
	for _, command := range config.Commands {
		if command.Alias == "" {
			continue
		}
		description := command.Description
		if description == "" {
			description = command.Name
		}
		if len(command.Alias) > leftColumnAliasLength {
			spaces := strings.Repeat(" ", leftColumnAliasLength+1)
			fmt.Printf("  %s\n  %s%s\n", command.Alias, spaces, description)
		} else {
			spaces := strings.Repeat(" ", leftColumnAliasLength+1-len(command.Alias))
			fmt.Printf("  %s%s%s\n", command.Alias, spaces, description)
		}
	}
}
