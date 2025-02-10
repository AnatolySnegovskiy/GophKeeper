package main

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
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
  port: 5555
redis:
  host: localhost
  port: 6379
  password: secret
db:
  dsn: postgres://user:password@localhost/dbname
`
	gormPostgres = func(dsn string) (*gorm.DB, error) {
		return nil, nil
	}

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	a := &app{}
	a.makeApp()
}
