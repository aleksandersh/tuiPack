package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aleksandersh/tuiPack/config"
	"github.com/aleksandersh/tuiPack/executor"
	"github.com/aleksandersh/tuiPack/tui"
	"github.com/alexflint/go-arg"
)

const (
	leftColumnAliasLength = 30
)

type args struct {
	Config  string `arg:"-c,--config" default:"./tuiPackConfig.yml" help:"path to config file"`
	Script  string `arg:"-s,--script" help:"run script by the alias"`
	Aliases bool   `arg:"-a,--aliases" help:"print list of aliases for available scripts"`
}

func main() {
	args := &args{}
	arg.MustParse(args)

	if args.Config == "" {
		args.Config = "./tuiPackConfig.yml"
	}
	config, err := config.ReadConfigFromYamlFile(args.Config)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if args.Aliases {
		printAliases(args, config)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if args.Script != "" {
		executeScript(ctx, args.Script, config)
		return
	}

	if err = tui.RunApp(ctx, config); err != nil {
		log.Fatalf("failed to run tui application: %v", err)
	}
}

func printAliases(args *args, config *config.Pack) {
	fmt.Printf("Usage: tuiPack --script ALIAS\n\n")
	fmt.Printf("Aliases for available scripts:\n")
	for _, command := range config.Commands {
		description := command.Description
		if description == "" {
			description = command.Name
		}
		if len(command.Alias) > leftColumnAliasLength {
			spaces := strings.Repeat(" ", leftColumnAliasLength+1)
			fmt.Printf("  %s\n  %s%s\n", command.Alias, spaces, description)
		} else {
			spaces := strings.Repeat(" ", leftColumnAliasLength+1-len(command.Alias))
			fmt.Printf("  %s%s%s\n", command.Alias, spaces, description)
		}
	}
}

func executeScript(ctx context.Context, alias string, config *config.Pack) {
	for _, command := range config.Commands {
		if command.Alias == alias {
			executor.ExecuteCommand(ctx, command.Args)
			return
		}
	}
	log.Fatalf("failed to find command by alias: %s", alias)
}
