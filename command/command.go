package command

import (
	"context"

	"github.com/aleksandersh/tuiPack/application"
)

type CommandFactory interface {
	CreateComand(properties *Properties) *Command
}

type Command interface {
	GetProperties() *Properties
	Execute(ctx context.Context, app *application.Application)
}

type Properties struct {
	Name        string
	Description string
	Alias       string
}

func NewProperties(Name string, Description string, Alias string) *Properties {
	return &Properties{
		Name:        Name,
		Description: Description,
		Alias:       Alias,
	}
}
