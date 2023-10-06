package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mattn/go-shellwords"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Script      string `yaml:"script"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing config argument")
	}
	configFile := os.Args[1]
	config, err := readConfig(configFile)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	app := tview.NewApplication()
	contentView, err := createContentView(app, config)
	if err != nil {
		log.Fatalf("failed to create content view: %v", err)
	}

	app.SetRoot(contentView, true).SetFocus(contentView)
	if err := app.Run(); err != nil {
		log.Fatalf("failed to run tui application: %v", err)
	}
	app.Stop()
}

func readConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in os.ReadFile: %w", err)
	}
	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error in yaml.Unmarshal: %w", err)
	}
	return &config, nil
}

func createContentView(app *tview.Application, config *Config) (tview.Primitive, error) {
	parser := shellwords.NewParser()
	parser.ParseEnv = true

	commandsView := tview.NewList()

	commandsView.SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetWrapAround(false).
		SetTitle(config.Name).
		SetBorder(true)

	for _, command := range config.Commands {
		err := addCommand(app, parser, commandsView, &command)
		if err != nil {
			return nil, err
		}
	}

	commandsView.Focus(func(p tview.Primitive) {})

	return commandsView, nil
}

func addCommand(app *tview.Application, parser *shellwords.Parser, listView *tview.List, command *Command) error {
	args, err := parser.Parse(command.Script)
	if err != nil {
		return fmt.Errorf("error in parser.Parse: %w", err)
	}

	listView.AddItem(command.Name, command.Description, 0, func() {
		execute(app, args)
	})
	return nil
}

func execute(app *tview.Application, args []string) {
	name := args[0]
	args = args[1:]
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	app.Stop()
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Run: %v", err)
		}
	}
}
