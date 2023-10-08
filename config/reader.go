package config

import (
	"errors"
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

var (
	errorFailedToParseEnvironment = errors.New("failed to parse environment")
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
	parsedCommands := make([]Command, 0, len(pack.Commands))
	for _, command := range pack.Commands {
		if command.Script == "" {
			return fmt.Errorf("missing command script: %v", command)
		}

		envEntries, err := parseEnvEntries(command.Environment)
		if err != nil {
			return err
		}

		parser.Getenv = func(env string) string {
			if env == ENV_TUI_PACK_EXECUTION_DIR {
				return executionDir
			}
			if env == ENV_TUI_PACK_CONFIG_DIR {
				return configDir
			}
			value := envEntries[env]
			if value != "" {
				return value
			}
			return os.Getenv(env)
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
		envExecutionDir := fmt.Sprintf("%s=%s", ENV_TUI_PACK_EXECUTION_DIR, executionDir)
		envConfigDir := fmt.Sprintf("%s=%s", ENV_TUI_PACK_CONFIG_DIR, configDir)
		environment := append(command.Environment, envExecutionDir, envConfigDir)
		parsedCommands = append(parsedCommands, Command{
			Name:        name,
			Description: command.Description,
			Script:      command.Script,
			Alias:       command.Alias,
			Environment: environment,
			Args:        args,
		})
	}
	pack.Commands = parsedCommands
	return nil
}

func parseEnvEntries(env []string) (map[string]string, error) {
	commandEnv := make(map[string]string, len(env))
	for _, entry := range env {
		key, value, err := parseEnvEntry(entry)
		if err != nil {
			return map[string]string{}, fmt.Errorf("failed to parse environment: env=%s, error=%w", entry, err)
		}
		commandEnv[key] = value
	}
	return commandEnv, nil
}

func parseEnvEntry(entry string) (string, string, error) {
	lastIndex := len(entry) - 1
	for i, s := range entry {
		if s == '=' {
			key := entry[:i]
			if i == lastIndex {
				return key, "", nil
			}
			value := entry[i+1:]
			return key, value, nil
		}
	}
	return "", "", errorFailedToParseEnvironment
}
