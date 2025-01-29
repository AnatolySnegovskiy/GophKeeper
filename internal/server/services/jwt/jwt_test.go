package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	jwtInstance := NewJwt()

	// Создаем токен
	token, err := jwtInstance.CreateToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем, что токен не истек
	err = jwtInstance.CheckToken(token)
	assert.NoError(t, err)
}

func TestCheckTokenExpired(t *testing.T) {
	jwtInstance := NewJwt()

	// Создаем токен с истекшим сроком действия
	jwtInstance.ExpiredAt = -time.Minute * 1
	token, err := jwtInstance.CreateToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем, что токен истек
	err = jwtInstance.CheckToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, "Token is expired")
}

func TestCheckTokenInvalid(t *testing.T) {
	jwtInstance := NewJwt()

	// Создаем некорректный токен
	invalidToken := "invalid.token.string"

	// Проверяем, что токен некорректен
	err := jwtInstance.CheckToken(invalidToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestGetExpiredAt(t *testing.T) {
	jwtInstance := NewJwt()

	// Проверяем, что метод GetExpiredAt возвращает правильное значение
	expiredAt := jwtInstance.GetExpiredAt()
	assert.Equal(t, time.Minute*15, expiredAt)
}
