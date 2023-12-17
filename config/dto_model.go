package config

type packDto struct {
	Name        string       `yaml:"name"`
	Version     string       `yaml:"version"`
	Environment []string     `yaml:"environment"`
	Commands    []commandDto `yaml:"commands"`
}

type commandDto struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Script      string   `yaml:"script"`
	Alias       string   `yaml:"alias"`
	Environment []string `yaml:"environment"`
}
