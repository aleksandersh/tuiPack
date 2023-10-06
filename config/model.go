package config

type Pack struct {
	Name     string    `yaml:"name"`
	Version  string    `yaml:"version"`
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Script      string `yaml:"script"`
	Alias       string `yaml:"alias"`
	Args        []string
}
