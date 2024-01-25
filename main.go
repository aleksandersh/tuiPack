package main

import (
	"context"
	"log"

	"github.com/aleksandersh/tuiPack/cli"
	"github.com/aleksandersh/tuiPack/loader"
	"github.com/aleksandersh/tuiPack/pack"
	"github.com/aleksandersh/tuiPack/pack/command"
	"github.com/aleksandersh/tuiPack/pack/command/inlined"
	"github.com/aleksandersh/tuiPack/pack/command/script"
)

func main() {
	args := cli.GetArgs()

	parsers := map[string]pack.Parser{
		command.CommandTypeScript:      script.NewParser(),
		command.CommandTypeInlinedPack: inlined.NewParser(),
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
