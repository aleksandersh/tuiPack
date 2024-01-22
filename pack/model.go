package pack

import "github.com/aleksandersh/tuiPack/command"

type Pack struct {
	Name            string
	CommandEntities []command.CommandEntity
}

func New(Name string, Commands []command.CommandEntity) *Pack {
	return &Pack{
		Name:            Name,
		CommandEntities: Commands,
	}
}
