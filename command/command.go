package command

import (
	"context"

	"github.com/aleksandersh/tuiPack/application"
)

type CommandFactory interface {
	CreateComand(properties *Properties) *CommandEntity
}

type CommandEntity struct {
	Command    Command
	Properties *Properties
}

type Properties struct {
	Name        string
	Description string
	Alias       string
}

type Command interface {
	Execute(ctx context.Context, app *application.Application, props *Properties)
}

func NewProperties(Name string, Description string, Alias string) *Properties {
	return &Properties{
		Name:        Name,
		Description: Description,
		Alias:       Alias,
	}
}

func NewEntity(properties *Properties, command Command) *CommandEntity {
	return &CommandEntity{Command: command, Properties: properties}
}
