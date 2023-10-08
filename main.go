package main

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/app"
	"github.com/aleksandersh/tuiPack/app/command"
	"github.com/aleksandersh/tuiPack/app/config"
)

func main() {
	args := app.GetArgs()

	config, err := config.ReadConfigFromYamlFile(args.Config)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if args.Aliases {
		command.PrintAliases(args, config)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if args.Script != "" {
		command.ExecuteScript(ctx, config, args.Script)
		return
	}

	command.RunTuiApp(ctx, config)
}
