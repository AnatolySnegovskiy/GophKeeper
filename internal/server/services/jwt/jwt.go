package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	key       []byte
	ExpiredAt time.Duration
}

func NewJwt() *Jwt {
	return &Jwt{
		key:       []byte("e45trjdmjkei54kmrdjkjei54eew"),
		ExpiredAt: time.Minute * 15,
	}
}

func (j *Jwt) GetExpiredAt() time.Duration {
	return j.ExpiredAt
}

func (j *Jwt) CreateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	claims["exp"] = time.Now().Add(j.ExpiredAt).Unix()
	tokenString, err := token.SignedString(j.key)
	return tokenString, err
}

func (j *Jwt) CheckToken(signedToken string) error {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.key, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return fmt.Errorf("token expired")
		}
		return nil
	}

	return fmt.Errorf("invalid token")
}
