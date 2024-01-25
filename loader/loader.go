package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/aleksandersh/tuiPack/pack"
	"github.com/aleksandersh/tuiPack/pack/command"
)

const (
	envFileDir      = "COMMAND_PACK_DIR"
	envExecutionDir = "COMMAND_PACK_EXECUTION_DIR"
)

type Loader struct {
	parsers map[string]pack.Parser
}

type loaderContext struct {
	env    map[string]string
	loader *Loader
	pack   *pack.Pack
}

func New(parsers map[string]pack.Parser) *Loader {
	return &Loader{parsers: parsers}
}

func (context *loaderContext) UnmarshalTOML(data interface{}) error {
	root, success := data.(map[string]interface{})
	if !success {
		return fmt.Errorf("failed to parse pack")
	}

	packEnvData, success := root["env"]
	if success {
		packEnv, success := packEnvData.(map[string]interface{})
		if !success {
			return fmt.Errorf("failed to parse pack environment")
		}
		for key, value := range packEnv {
			str, success := value.(string)
			if !success {
				return fmt.Errorf("failed to parse pack environment: key=%s", key)
			}
			context.env[key] = str
		}
	}

	commands, success := root["commands"].([]map[string]interface{})
	if !success {
		return fmt.Errorf("failed to parse commands")
	}
	idx := 0
	commandEntities := make([]command.CommandEntity, 0, len(commands))
	for _, cmd := range commands {
		entities, err := context.parseCommand(cmd, idx)
		if err != nil {
			return fmt.Errorf("failed to parse commands: %w", err)
		}
		commandEntities = append(commandEntities, entities...)
		idx++
	}

	packName := root["name"].(string)
	context.pack = &pack.Pack{
		Name:            packName,
		CommandEntities: commandEntities,
	}
	return nil
}

func (context *loaderContext) parseCommand(data map[string]interface{}, idx int) ([]command.CommandEntity, error) {
	cmdType, success := data["type"].(string)
	if !success {
		cmdType = command.CommandTypeScript
	}
	parser, success := context.loader.parsers[cmdType]
	if !success {
		return nil, fmt.Errorf("missing parser: type=%s, index=%d", cmdType, idx)
	}
	cmdName, success := data["name"].(string)
	if !success {
		return nil, fmt.Errorf("failed to parse command name: index=%d", idx)
	}
	var cmdDescription string
	cmdDescriptionData, success := data["description"]
	if success {
		cmdDescription, success = cmdDescriptionData.(string)
		if !success {
			return nil, fmt.Errorf("failed to parse command description: index=%d", idx)
		}
	}
	var cmdAlias string
	cmdAliasData, success := data["alias"]
	if success {
		cmdAlias, success = cmdAliasData.(string)
		if !success {
			return nil, fmt.Errorf("failed to parse command alias: index=%d", idx)
		}
	}
	cmdEnv := make(map[string]string, len(context.env))
	for key, value := range context.env {
		cmdEnv[key] = value
	}
	cmdEnvData, success := data["env"]
	if success {
		cmdEnvMap, success := cmdEnvData.(map[string]interface{})
		if !success {
			return nil, fmt.Errorf("failed to parse command environment: index=%d", idx)
		}
		for key, value := range cmdEnvMap {
			str, success := value.(string)
			if !success {
				return nil, fmt.Errorf("failed to parse command environment: index=%d, key=%s", idx, key)
			}
			cmdEnv[key] = str
		}
	}

	properties := command.NewProperties(cmdName, cmdDescription, cmdAlias, cmdEnv)
	entities, err := parser.Parse(data, properties, context)
	if err != nil {
		return nil, fmt.Errorf("error in parser.Parse: %w", err)
	}
	return entities, nil
}

func (loader *Loader) Load(path string) (*pack.Pack, error) {
	configDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("error in filepath.Abs: %w", err)
	}
	executionDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error in os.Getwd: %w", err)
	}

	env := make(map[string]string, 2)
	env[envFileDir] = configDir
	env[envExecutionDir] = executionDir

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in os.ReadFile: %w", err)
	}

	newContext := &loaderContext{env: env, loader: loader}
	err = toml.Unmarshal(file, newContext)
	if err != nil {
		return nil, fmt.Errorf("error in toml.Unmarshal: %w", err)
	}
	return newContext.pack, nil
}

func (context *loaderContext) Load(path string) (*pack.Pack, error) {
	configDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("error in filepath.Abs: %w", err)
	}

	env := make(map[string]string, len(context.env))
	for key, value := range context.env {
		env[key] = value
	}
	env[envFileDir] = configDir

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in os.ReadFile: %w", err)
	}

	newContext := &loaderContext{env: env, loader: context.loader}
	err = toml.Unmarshal(file, newContext)
	if err != nil {
		return nil, fmt.Errorf("error in toml.Unmarshal: %w", err)
	}
	return newContext.pack, nil
}
