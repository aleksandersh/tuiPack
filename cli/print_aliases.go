package cli

import (
	"fmt"
	"strings"

	"github.com/aleksandersh/tuiPack/config"
)

const (
	leftColumnAliasLength = 30
)

func PrintAliases(args *Args, config *config.Pack) {
	fmt.Printf("Usage: tuiPack --script ALIAS\n\n")
	fmt.Printf("Aliases for available scripts:\n")
	for _, entity := range config.CommandEntities {
		properties := entity.Properties
		if properties.Alias == "" {
			continue
		}
		description := properties.Description
		if description == "" {
			description = properties.Name
		}
		if len(properties.Alias) > leftColumnAliasLength {
			spaces := strings.Repeat(" ", leftColumnAliasLength+1)
			fmt.Printf("  %s\n  %s%s\n", properties.Alias, spaces, description)
		} else {
			spaces := strings.Repeat(" ", leftColumnAliasLength+1-len(properties.Alias))
			fmt.Printf("  %s%s%s\n", properties.Alias, spaces, description)
		}
	}
}
