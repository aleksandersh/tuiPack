package command

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/app/config"
	"github.com/aleksandersh/tuiPack/tui"
)

func RunTuiApp(ctx context.Context, config *config.Pack) {
	if err := tui.RunApp(ctx, config); err != nil {
		log.Fatalf("failed to run tui application: %v", err)
	}
}
