package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/aleksandersh/tuiPack/command"
	"github.com/aleksandersh/tuiPack/pack"
)

const (
	envFileDir      = "COMMAND_PACK_DIR"
	envExecutionDir = "COMMAND_PACK_EXECUTION_DIR"
)

type Loader struct {
	parsers map[string]command.Parser
}

type packDto struct {
	fileDir      string
	executionDir string
	loader       *Loader
	pack         *pack.Pack
}

func New(parsers map[string]command.Parser) *Loader {
	return &Loader{parsers: parsers}
}

func (dto *packDto) UnmarshalTOML(data interface{}) error {
	root, success := data.(map[string]interface{})
	if !success {
		return fmt.Errorf("failed to parse pack")
	}

	strEnv := make(map[string]string)
	strEnv[envFileDir] = dto.fileDir
	strEnv[envExecutionDir] = dto.executionDir

	packEnv, success := root["env"].(map[string]interface{})
	if !success {
		return fmt.Errorf("failed to parse pack environment")
	}
	for key, value := range packEnv {
		str, success := value.(string)
		if !success {
			return fmt.Errorf("failed to parse pack environment: key=%s", key)
		}
		strEnv[key] = str
	}

	commands, success := root["commands"].([]map[string]interface{})
	if !success {
		return fmt.Errorf("failed to parse commands")
	}
	idx := 0
	commandEntities := make([]command.CommandEntity, 0, len(commands))
	for _, cmd := range commands {
		entity, err := dto.parseCommand(cmd, strEnv, idx)
		if err != nil {
			return fmt.Errorf("failed to parse commands: %w", err)
		}
		commandEntities = append(commandEntities, *entity)
		idx++
	}

	packName := root["name"].(string)
	dto.pack = &pack.Pack{
		Name:            packName,
		CommandEntities: commandEntities,
	}
	return nil
}

func (dto *packDto) parseCommand(data map[string]interface{}, env map[string]string, idx int) (*command.CommandEntity, error) {
	cmdType, success := data["type"].(string)
	if !success {
		return nil, fmt.Errorf("failed to parse command type: index=%d", idx)
	}
	parser, success := dto.loader.parsers[cmdType]
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
	cmdEnv := make(map[string]string, 0)
	for key, value := range env {
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
	cmd, err := parser.Parse(data, properties)
	if err != nil {
		return nil, fmt.Errorf("error in parser.Parse: %w", err)
	}
	return cmd, nil
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

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error in os.ReadFile: %w", err)
	}

	dto := &packDto{fileDir: configDir, executionDir: executionDir, loader: loader}
	err = toml.Unmarshal(file, dto)
	if err != nil {
		return nil, fmt.Errorf("error in toml.Unmarshal: %w", err)
	}
	return dto.pack, nil
}
