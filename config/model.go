package config

import "github.com/aleksandersh/tuiPack/command"

type Pack struct {
	Name            string
	Version         string
	CommandEntities []command.CommandEntity
}

func New(Name string, Version string, Commands []command.CommandEntity) *Pack {
	return &Pack{
		Name:            Name,
		Version:         Version,
		CommandEntities: Commands,
	}
}
