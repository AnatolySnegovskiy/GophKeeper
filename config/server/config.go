package server

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	}
	DB struct {
		DSN string `yaml:"dsn"`
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

	configFile := filepath.Join(cwd, "config.yaml")
	file, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, c)
	return err
}
