package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paulozy/costurai/configs"
)

type GenerateTokenInput struct {
	Issuer  string
	Subject string
}

func GenerateToken(data GenerateTokenInput) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": data.Issuer,
		"sub": data.Subject,
		"exp": time.Now().Add(time.Hour * time.Duration(getJWTExpiration())).Unix(),
	})

	tokenString, err := token.SignedString([]byte(getJWTSecretKey()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(getJWTSecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims
	}

	return token, nil
}

func getJWTSecretKey() string {
	config, err := configs.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	return config.JWTSecret
}

func getJWTExpiration() int64 {
	config, err := configs.LoadConfig("../")
	if err != nil {
		panic(err)
	}

	return config.JWTExpiresIn
}
