package config

import (
	"fmt"
	"os"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

func ReadConfigFromYamlFile(path string) (*Pack, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in os.ReadFile: %w", err)
	}
	pack := &Pack{}
	if err = yaml.Unmarshal(file, pack); err != nil {
		return nil, fmt.Errorf("error in yaml.Unmarshal: %w", err)
	}

	if err = parseCommandsArgs(pack); err != nil {
		return nil, err
	}

	return pack, nil
}

func parseCommandsArgs(pack *Pack) error {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	parsedCommands := make([]Command, 0, len(pack.Commands))
	for _, command := range pack.Commands {
		args, err := parser.Parse(command.Script)
		if err != nil {
			return fmt.Errorf("error in parser.Parse: %w", err)
		}
		parsedCommands = append(parsedCommands, Command{
			Name:        command.Name,
			Description: command.Description,
			Script:      command.Script,
			Args:        args,
		})
	}
	pack.Commands = parsedCommands
	return nil
}
