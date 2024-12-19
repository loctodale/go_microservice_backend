package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"log"
	"net/http"
	"time"
)

var m = make(map[string][]string)

type CustomClaims struct {
	UserID string   `json:"sub"`
	Roles  []string `json:"roles"`
	Scope  []string `json:"scp"`
	jwt.RegisteredClaims
}
type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTService() *JWTService {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		global.Logger.Fatal("Fail to generate private key", zap.Error(err))
	}

	return &JWTService{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

func (s *JWTService) GenerateToken(userId string, role string) (string, error) {
	RoleSet("shop", []string{"write:shop", "read:shop"})
	claims := CustomClaims{
		UserID: userId,
		Roles:  []string{role},
		Scope:  m[role],
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://localhost:4455/",
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"https://localhost:4455/"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		global.Logger.Error("Fail to generate token", zap.Error(err))
	}

	return tokenString, nil
}

func (s *JWTService) GenerateJWKS() *jose.JSONWebKeySet {
	key := jose.JSONWebKey{
		Key:       s.publicKey,
		KeyID:     "key-1", // Unique key identifier
		Algorithm: "RS256",
		Use:       "sig", // Sử dụng cho việc ký
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{key},
	}
}

func (s *JWTService) JWKSHandler(w http.ResponseWriter, r *http.Request) {
	jwks := s.GenerateJWKS()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

func RoleSet(role string, permissions []string) {
	m[role] = permissions
}
func (s *JWTService) GenerateRefreshToken(userID string) string {
	refreshToken := jwt.New(jwt.SigningMethodRS256)
	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix() // 7 days

	token, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		log.Fatalf("Failed to create refresh token: %v", err)
	}

	return token
}
