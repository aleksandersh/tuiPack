package command

import (
	"context"

	"github.com/aleksandersh/tuiPack/application"
)

const (
	CommandTypeScript      = "script"
	CommandTypeInlinedPack = "inlined_pack"
)

type Command interface {
	Execute(ctx context.Context, app *application.Application, props *Properties)
}

type CommandEntity struct {
	Command    Command
	Properties *Properties
}

type Properties struct {
	Name        string
	Description string
	Alias       string
	Environment map[string]string
}

func NewProperties(name string, description string, alias string, env map[string]string) *Properties {
	return &Properties{
		Name:        name,
		Description: description,
		Alias:       alias,
		Environment: env,
	}
}

func NewEntity(properties *Properties, command Command) *CommandEntity {
	return &CommandEntity{Command: command, Properties: properties}
}
