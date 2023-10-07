package main

import (
	"context"
	"log"
	"os"

	"github.com/aleksandersh/tuiPack/config"
	"github.com/aleksandersh/tuiPack/tui"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing config argument")
	}
	configFile := os.Args[1]
	config, err := config.ReadConfigFromYamlFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = tui.RunApp(ctx, config); err != nil {
		log.Fatalf("failed to run tui application: %v", err)
	}
}
