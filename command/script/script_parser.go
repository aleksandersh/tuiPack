package script

import (
	"fmt"

	"github.com/aleksandersh/tuiPack/command"
)

type ScriptParser struct {
}

func NewParser() *ScriptParser {
	return &ScriptParser{}
}

func (parser *ScriptParser) Parse(data map[string]interface{}, properties *command.Properties) (*command.CommandEntity, error) {
	script, success := data["script"].(string)
	if !success {
		return nil, fmt.Errorf("failed to parse script")
	}
	cmd := newScript(properties, script)
	return command.NewEntity(properties, cmd), nil
}
