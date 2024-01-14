package cli

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/config"
)

func ExecuteScript(ctx context.Context, config *config.Pack, alias string) {
	for _, entity := range config.CommandEntities {
		if entity.Properties.Alias == alias {
			app := application.NewApplication(application.NewController(nil))
			entity.Command.Execute(ctx, app, entity.Properties)
			return
		}
	}
	log.Fatalf("failed to find command by alias: %s", alias)
}
