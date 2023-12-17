package main

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/cli"
	"github.com/aleksandersh/tuiPack/launcher"
)

func main() {
	args := cli.GetArgs()

	config, err := launcher.ReadConfigFromYamlFile(args.Config)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if args.Aliases {
		cli.PrintAliases(args, config)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if args.Script != "" {
		cli.ExecuteScript(ctx, config, args.Script)
		return
	}

	cli.RunTuiApp(ctx, config)
}
