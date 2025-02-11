package main

import (
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
	redisClient = func(opt *redis.Options) *redis.Client {
		mockRedis, mock := redismock.NewClientMock()
		mock.ExpectPing().SetVal("pong")
		return mockRedis
	}
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	a := &app{}
	a.makeApp()
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
	configFile := filepath.Join(configDir, "config-server.yaml")
	configContent := `
server
redis:
  host: localhost
  por22t: 637
  pass2word: secret
db:
  dsn: postgres://user:password@localhost/dbname
`
	gormPostgres = func(dsn string) (*gorm.DB, error) {
		return nil, nil
	}
	redisClient = func(opt *redis.Options) *redis.Client {
		mockRedis, mock := redismock.NewClientMock()
		mock.ExpectPing().SetVal("pong")
		return mockRedis
	}
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	a := &app{}
	go a.makeApp()
	os.RemoveAll(configDir)
}
