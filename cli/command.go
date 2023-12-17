package cli

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/config"
)

func ExecuteScript(ctx context.Context, config *config.Pack, alias string) {
	for _, command := range config.Commands {
		properties := command.GetProperties()
		if properties.Alias == alias {
			app := application.NewApplication(application.NewController(nil))
			command.Execute(ctx, app)
			return
		}
	}
	log.Fatalf("failed to find command by alias: %s", alias)
}
