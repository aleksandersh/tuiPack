package config

type Pack struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Environment []string  `yaml:"environment"`
	Commands    []Command `yaml:"commands"`
}

type Command struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Script      string   `yaml:"script"`
	Alias       string   `yaml:"alias"`
	Environment []string `yaml:"environment"`
	Args        []string
}
