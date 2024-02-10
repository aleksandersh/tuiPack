package inlined

import (
	"fmt"
	"os"

	"github.com/aleksandersh/tuiPack/pack"
	"github.com/aleksandersh/tuiPack/pack/command"
	"github.com/mattn/go-shellwords"
)

type inlinedPackParser struct {
}

func NewParser() *inlinedPackParser {
	return &inlinedPackParser{}
}

func (parser *inlinedPackParser) Parse(data map[string]interface{}, properties *command.Properties, loader pack.Loader) ([]command.CommandEntity, error) {
	path, success := data["path"].(string)
	if !success {
		return nil, fmt.Errorf("failed to parse script")
	}

	path, err := parser.resolvePathWithEnv(path, properties)
	if err != nil {
		return nil, err
	}

	namePrefix := ""
	namePrefixData, success := data["name_prefix"]
	if success {
		namePrefix, success = namePrefixData.(string)
		if !success {
			return nil, fmt.Errorf("failed to parse name_prefix")
		}
	}

	aliasPrefix := ""
	aliasPrefixData, success := data["alias_prefix"]
	if success {
		aliasPrefix, success = aliasPrefixData.(string)
		if !success {
			return nil, fmt.Errorf("failed to parse alias_prefix")
		}
	}

	pack, err := loader.Load(path)
	if err != nil {
		return nil, fmt.Errorf("error in loader.Load: %w", err)
	}
	entities := pack.CommandEntities
	if namePrefix != "" || aliasPrefix != "" {
		for _, entity := range entities {
			entity.Properties.Name = namePrefix + entity.Properties.Name
			if entity.Properties.Alias != "" {
				entity.Properties.Alias = aliasPrefix + entity.Properties.Alias
			}
		}
	}
	return entities, nil
}

func (cmd *inlinedPackParser) resolvePathWithEnv(path string, props *command.Properties) (string, error) {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	parser.Getenv = func(env string) string {
		if value, contains := props.Environment[env]; contains {
			return value
		}
		return os.Getenv(env)
	}
	args, err := parser.Parse(path)
	if err != nil {
		return "", fmt.Errorf("error in parser.Parse: %w", err)
	}
	if len(args) != 1 {
		return "", fmt.Errorf("failed to process pack path: %s", path)
	}
	return args[0], nil
}
