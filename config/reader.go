package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

const (
	ENV_TUI_PACK_CONFIG_DIR    = "TUI_PACK_CONFIG_DIR"
	ENV_TUI_PACK_EXECUTION_DIR = "TUI_PACK_EXECUTION_DIR"
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

	configDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("error in filepath.Abs: %w", err)
	}

	executionDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error in os.Getwd: %w", err)
	}

	if err = parseCommandsArgs(configDir, executionDir, pack); err != nil {
		return nil, err
	}

	return pack, nil
}

func parseCommandsArgs(configDir string, executionDir string, pack *Pack) error {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	parser.Getenv = func(env string) string {
		if env == ENV_TUI_PACK_EXECUTION_DIR {
			return executionDir
		}
		if env == ENV_TUI_PACK_CONFIG_DIR {
			return configDir
		}
		return os.Getenv(env)
	}
	parsedCommands := make([]Command, 0, len(pack.Commands))
	for _, command := range pack.Commands {
		if command.Script == "" {
			return fmt.Errorf("missing command script: %v", command)
		}
		args, err := parser.Parse(command.Script)
		if err != nil {
			return fmt.Errorf("error in parser.Parse: %w", err)
		}
		name := command.Name
		if name == "" {
			if command.Alias != "" {
				name = command.Alias
			} else {
				name = command.Script
			}
		}
		parsedCommands = append(parsedCommands, Command{
			Name:        name,
			Description: command.Description,
			Script:      command.Script,
			Alias:       command.Alias,
			Args:        args,
		})
	}
	pack.Commands = parsedCommands
	return nil
}
