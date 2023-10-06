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
	app := tview.NewApplication()
	contentView := createContentView(app)
	app.SetRoot(contentView, true).SetFocus(contentView)
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
	app.Stop()
}

func createContentView(app *tview.Application) tview.Primitive {
	file, err := os.ReadFile("./tuiPackConfig.yml")
	if err != nil {
		log.Fatal(err)
	}
	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	parser := shellwords.NewParser()
	parser.ParseEnv = true

	commandsView := tview.NewList()

	commandsView.SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetWrapAround(false).
		SetTitle(config.Name).
		SetBorder(true)

	for _, command := range config.Commands {
		addCommand(app, parser, commandsView, &command)
	}

	commandsView.Focus(func(p tview.Primitive) {})

	return commandsView
}

func addCommand(app *tview.Application, parser *shellwords.Parser, listView *tview.List, command *Command) {
	args, err := parser.Parse(command.Script)
	if err != nil {
		log.Fatal(err)
	}

	listView.AddItem(command.Name, command.Description, 0, func() {
		execute(app, args)
	})
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
		log.Fatal(err)
	}
}
