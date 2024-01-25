package pack

import "github.com/aleksandersh/tuiPack/pack/command"

type Loader interface {
	Load(path string) (*Pack, error)
}

type Parser interface {
	Parse(data map[string]interface{}, properties *command.Properties, loader Loader) ([]command.CommandEntity, error)
}

type Pack struct {
	Name            string
	CommandEntities []command.CommandEntity
}

func NewPack(Name string, Commands []command.CommandEntity) *Pack {
	return &Pack{
		Name:            Name,
		CommandEntities: Commands,
	}
}
