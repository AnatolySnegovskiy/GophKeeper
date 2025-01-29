package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// Создаем директорию для конфигурационного файла
	configDir := filepath.Join(cwd, "configs")
	err = os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Создаем временный файл конфигурации для тестирования
	configFile := filepath.Join(configDir, "config-server.yaml")
	configContent := `
server:
  host: localhost
  port: 8080
redis:
  host: localhost
  port: 6379
  password: secret
db:
  dsn: postgres://user:password@localhost:5432/dbname
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Тестируем функцию NewConfig
	config, err := NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "localhost", config.Redis.Host)
	assert.Equal(t, "6379", config.Redis.Port)
	assert.Equal(t, "secret", config.Redis.Password)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", config.DB.DSN)

	// Удаляем созданный файл и директорию после теста
	defer os.Remove(configFile)
	defer os.Remove(configDir)
}

func TestParseYaml(t *testing.T) {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// Создаем директорию для конфигурационного файла
	configDir := filepath.Join(cwd, "configs")
	err = os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Создаем временный файл конфигурации для тестирования
	configFile := filepath.Join(configDir, "config-server.yaml")
	configContent := `
server:
  host: localhost
  port: 8080
redis:
  host: localhost
  port: 6379
  password: secret
db:
  dsn: postgres://user:password@localhost:5432/dbname
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Создаем экземпляр Config
	config := &Config{}

	// Тестируем метод parseYaml
	err = config.parseYaml()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "localhost", config.Redis.Host)
	assert.Equal(t, "6379", config.Redis.Port)
	assert.Equal(t, "secret", config.Redis.Password)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", config.DB.DSN)

	// Удаляем созданный файл и директорию после теста
	defer os.Remove(configFile)
	defer os.Remove(configDir)
}

func TestParseYamlError(t *testing.T) {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// Создаем директорию для конфигурационного файла
	configDir := filepath.Join(cwd, "configs")
	err = os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Создаем временный файл конфигурации с некорректным содержимым для тестирования
	configFile := filepath.Join(configDir, "config-server.yaml")
	configContent := `
invalid:
  host: localhost
  pot 8080
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Создаем экземпляр Config
	config := &Config{}

	// Тестируем метод parseYaml с некорректным содержимым
	err = config.parseYaml()
	assert.Error(t, err)

	// Удаляем созданный файл и директорию после теста
	defer os.Remove(configFile)
	defer os.Remove(configDir)
}
