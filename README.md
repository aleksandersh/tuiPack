# TUI commands pack

A tool for organizing frequently used commands within the terminal user interface

You can describe a set of commands in a single YAML file and easily access them from the terminal interface

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
name: My test commands pack
commands:
  - name: list files
    script: ls -l "$TUI_PACK_EXECUTION_DIR"
    alias: list_files
    description: list files in the current directory
  - name: sleep
    environment:
      - TIME_SECONDS=3
    script: sleep "$TIME_SECONDS"
    alias: sleep
    description: sleep for 3 seconds
```
