package cli

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/pack"
	"github.com/aleksandersh/tuiPack/tui"
)

func RunTuiApp(ctx context.Context, pack *pack.Pack) {
	if err := tui.RunApp(ctx, pack); err != nil {
		log.Fatalf("failed to run tui application: %v", err)
	}
}
