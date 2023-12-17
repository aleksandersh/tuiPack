package script

import (
	"github.com/aleksandersh/tuiPack/command"
)

type ScriptFactory struct {
}

func (sf ScriptFactory) CreateComand(properties *command.Properties, args []string, env []string) command.Command {
	return newScript(properties, args, env)
}
