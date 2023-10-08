package command

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/app/config"
	"github.com/aleksandersh/tuiPack/executor"
)

func ExecuteScript(ctx context.Context, config *config.Pack, alias string) {
	for _, command := range config.Commands {
		if command.Alias == alias {
			executor.ExecuteCommand(ctx, command.Args, command.Environment)
			return
		}
	}
	log.Fatalf("failed to find command by alias: %s", alias)
}
