package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aleksandersh/tuiPack/command"
	"github.com/aleksandersh/tuiPack/command/script"
	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

const (
	envConfigDir    = "TUI_PACK_CONFIG_DIR"
	envExecutionDir = "TUI_PACK_EXECUTION_DIR"
)

var (
	errorFailedToParseEnvironment = errors.New("failed to parse environment")
)

func ReadConfigFromYamlFile(path string) (*Pack, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in os.ReadFile: %w", err)
	}
	pack := &packDto{}
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

	commands, err := parseCommands(configDir, executionDir, pack)
	if err != nil {
		return nil, err
	}

	return New(pack.Name, pack.Version, commands), nil
}

func parseCommands(configDir string, executionDir string, pack *packDto) ([]command.CommandEntity, error) {
	envExecutionDirEntry := fmt.Sprintf("%s=%s", envExecutionDir, executionDir)
	envConfigDirEntry := fmt.Sprintf("%s=%s", envConfigDir, configDir)
	packEnvEntries, err := parseEnvEntries(pack.Environment)
	if err != nil {
		return nil, err
	}
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	scriptFactory := script.ScriptFactory{}
	parsedCommands := make([]command.CommandEntity, 0, len(pack.Commands))
	for _, cmd := range pack.Commands {
		if cmd.Script == "" {
			return nil, fmt.Errorf("missing command script: %v", cmd)
		}

		commandEnvEntries, err := parseEnvEntries(cmd.Environment)
		if err != nil {
			return nil, err
		}

		parser.Getenv = func(env string) string {
			if env == envExecutionDir {
				return executionDir
			}
			if env == envConfigDir {
				return configDir
			}
			if value, contains := commandEnvEntries[env]; contains {
				return value
			}
			if value, contains := packEnvEntries[env]; contains {
				return value
			}
			return os.Getenv(env)
		}
		args, err := parser.Parse(cmd.Script)
		if err != nil {
			return nil, fmt.Errorf("error in parser.Parse: %w", err)
		}

		name := cmd.Name
		if name == "" {
			if cmd.Alias != "" {
				name = cmd.Alias
			} else {
				name = cmd.Script
			}
		}
		environment := append(append(pack.Environment, cmd.Environment...), envExecutionDirEntry, envConfigDirEntry)
		properties := command.NewProperties(name, cmd.Description, cmd.Alias)
		script := scriptFactory.CreateComand(properties, args, environment)
		parsedCommands = append(parsedCommands, *script)
	}
	return parsedCommands, nil
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
