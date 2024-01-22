package main

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/cli"
	"github.com/aleksandersh/tuiPack/command"
	"github.com/aleksandersh/tuiPack/command/script"
	"github.com/aleksandersh/tuiPack/pack/loader"
)

func main() {
	args := cli.GetArgs()

	parsers := map[string]command.Parser{
		"script": script.NewParser(),
	}
	loader := loader.New(parsers)
	pack, err := loader.Load(args.Config)
	if err != nil {
		log.Fatalf("failed to load command pack: %v", err)
	}

	if args.Aliases {
		cli.PrintAliases(args, pack)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if args.Script != "" {
		cli.ExecuteScript(ctx, pack, args.Script)
		return
	}

	cli.RunTuiApp(ctx, pack)
}
