package config

import "github.com/aleksandersh/tuiPack/command"

type Pack struct {
	Name     string
	Version  string
	Commands []command.Command
}

func New(Name string, Version string, Commands []command.Command) *Pack {
	return &Pack{
		Name:     Name,
		Version:  Version,
		Commands: Commands,
	}
}
