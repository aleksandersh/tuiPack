# TUI commands pack

A tool for organizing frequently used commands within the terminal user interface.

You can describe a set of commands in a single YAML file and easily access them from the terminal interface.

Follow the example from the [examples catalog](./example).

![tuiPack example](./example/tuiPackExample.gif "Example")

## Usage

```bash
# ask for help
tuiPack --help
```

### Terminal user interface

```bash
# start the terminal user interface with commands from the specified config
tuiPack --config ./tuiPackConfig.yml
```

#### Hotkeys

`/` - to enter filtering mode

### Execute a command by an alias

```bash
# show available commands
tuiPack --config ./tuiPackConfig.yml --aliases
# execute "list_files" command
tuiPack --config ./tuiPackConfig.yml --script list_files
```

## YML config

##### additional environment variables

`TUI_PACK_CONFIG_DIR` - absolute path of the config directory  
`TUI_PACK_EXECUTION_DIR` - absolute path of the executuon directory  

##### example

```yml
name: command pack name
environment:
  - ENV_1=VALUE
commands:
  - name: command displayed name
    environment:
    - ENV_2=VALUE
    script: echo "$ENV_1" "$ENV_2"
    alias: command_alias_for_terminal
    description: description of the command
```
