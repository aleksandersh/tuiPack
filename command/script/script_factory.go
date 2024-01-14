package script

import "github.com/aleksandersh/tuiPack/command"

type ScriptFactory struct {
}

func (sf *ScriptFactory) CreateComand(properties *command.Properties, args []string, env []string) *command.CommandEntity {
	cmd := newScript(properties, args, env)
	return command.NewEntity(properties, cmd)
}
