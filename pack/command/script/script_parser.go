package script

import (
	"fmt"

	"github.com/aleksandersh/tuiPack/pack"
	"github.com/aleksandersh/tuiPack/pack/command"
)

type scriptParser struct {
}

func NewParser() *scriptParser {
	return &scriptParser{}
}

func (parser *scriptParser) Parse(data map[string]interface{}, properties *command.Properties, loader pack.Loader) ([]command.CommandEntity, error) {
	script, success := data["script"].(string)
	if !success {
		return nil, fmt.Errorf("failed to parse script")
	}
	cmd := newScript(script)
	return []command.CommandEntity{*command.NewEntity(properties, cmd)}, nil
}
