package cli

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/application"
	"github.com/aleksandersh/tuiPack/pack"
)

func ExecuteScript(ctx context.Context, pack *pack.Pack, alias string) {
	for _, entity := range pack.CommandEntities {
		if entity.Properties.Alias == alias {
			app := application.NewApplication(application.NewController(nil))
			entity.Command.Execute(ctx, app, entity.Properties)
			return
		}
	}
	log.Fatalf("failed to find a command by the alias: %s", alias)
}
