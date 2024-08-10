package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Jwt struct {
	key []byte
}

func NewJwt() *Jwt {
	return &Jwt{
		key: []byte("e45trjdmjkei54kmrdjkjei54eew"),
	}
}

func (j *Jwt) CreateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expired := time.Minute * 15
	claims["exp"] = time.Now().Add(expired).Unix()
	tokenString, err := token.SignedString(j.key)
	return tokenString, err
}

func (j *Jwt) Validate(signedToken string) error {
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
