package client

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}

func NewConfig() (Config, error) {
	config := Config{}
	err := config.parseYaml()
	return config, err
}

func (c *Config) parseYaml() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	configFile := filepath.Join(cwd, "configs/config-client.yaml")
	file, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, c)
	return err
}
