package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMakeApp(t *testing.T) {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// Создаем директорию для конфигурационного файла
	configDir := filepath.Join(cwd, "configs")
	err = os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Создаем временный файл конфигурации для тестирования
	configFile := filepath.Join(configDir, "config-client.yaml")
	configContent := `
server:
  host: localhost
  port: 8080
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	a := &app{}
	a.makeApp()
	a.stop()
	os.RemoveAll(configDir)
}

func TestMakeAppFail(t *testing.T) {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// Создаем директорию для конфигурационного файла
	configDir := filepath.Join(cwd, "configs")
	err = os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Создаем временный файл конфигурации для тестирования
	configFile := filepath.Join(configDir, "config-client.yaml")
	configContent := `
server:
  ho2st: localhost
  port: 8080
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	a := &app{}
	a.makeApp()
	os.RemoveAll(configDir)
}
