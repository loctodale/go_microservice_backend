package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go_microservice_backend_api/global"
	"time"
)

type PayloadClaims struct {
	jwt.RegisteredClaims
}

func GenTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.JWTSetting.APISecret))
}

func CreateToken(uuidToken string) (string, error) {
	//1. set expired time
	timeEx := global.Config.JWTSetting.JWTExpiration
	if timeEx == "" {
		timeEx = "1h"
	}
	exDuration, err := time.ParseDuration(timeEx)
	if err != nil {
		return "", err
	}
	expiredAt := time.Now().Add(exDuration)

	return GenTokenJWT(&PayloadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shopdevgo",
			Subject:   uuidToken,
		},
	})
}
