package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go_microservice_backend_api/global"
	"log"
	"os"
	"time"
)

type PayloadClaims struct {
	jwt.RegisteredClaims
}

func GenerateJWTSecret() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to base64 to get a URL-safe string
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func GenTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.JWTSetting.APISecret))
}

func GenTokenJWTPair(payload jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secret))
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

func loadPrivateKey() *rsa.PrivateKey {
	data, err := os.ReadFile("rsa_private.pem")
	if err != nil {
		log.Fatalf("Cannot read private key: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatalf("Invalid private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return key
}

// Generate access token (short-lived)
func GenerateAccessToken(userID string) (string, error) {
	privateKey := loadPrivateKey()
	timeEx, err := time.ParseDuration(global.Config.JWTSetting.JWTExpiration)
	if err != nil {
		return "Parse duration failed", err
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(timeEx).Unix(), // 15 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// Generate refresh token (long-lived)
