package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"log"
	"net/http"
	"os"
	"time"
)

var m = make(map[string][]string)

type Consumer struct {
	Username string `json:"username"`
	UserID   string `json:"custom_id"`
}
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
type CreateCredentialsResponse struct {
	Id           string `json:"id"`
	CreateAt     int    `json:"created_at"`
	RSAPublicKey string `json:"rsa_public_key"`
	Secret       string `json:"secret"`
	ConsumerId   string `json:"consumer.id"`
	Key          string `json:"key"`
	Alg          string `json:"algorithm"`
	Tags         string `json:"tags"`
}

func readPrivateKey() *rsa.PrivateKey {
	keyData, err := os.ReadFile("/etc/key/private_key.pem")
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	return key
}

func readPublicKey() *rsa.PublicKey {
	keyData, err := os.ReadFile("/etc/key/public_key.pem")
	if err != nil {
		log.Fatalf("Failed to read public key file: %v", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	return key
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
func GenerateRandomKID() (string, error) {
	// Generate 32 random bytes (256 bits)
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode the bytes to a URL-safe base64 string and return it
	kid := base64.URLEncoding.EncodeToString(randomBytes)
	return kid, nil
}

func publicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	// Chuyển đổi khóa công khai thành định dạng ASN.1 DER
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	// Mã hóa khóa công khai dưới dạng PEM
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	if publicKeyPEM == nil {
		return "", fmt.Errorf("failed to encode public key to PEM")
	}

	return string(publicKeyPEM), nil
}
func (s *JWTService) GenerateToken(userId string, name string, role string, key string) (accessToken string, err error) {
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

	token.Header["kid"] = key
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		global.Logger.Error("Fail to generate token", zap.Error(err))
	}

	return tokenString, nil

}

func (s *JWTService) GenerateTokenRegister(userId string, name string, role string) (accessToken string, credentialId string, err error) {
	key, err := GenerateRandomKID()
	if err != nil {
		global.Logger.Error("Fail to generate token", zap.Error(err))
		return "", "", err
	}
	tokenString, err := s.GenerateToken(userId, name, role, key)
	if err != nil {
		global.Logger.Error("Fail to generate token", zap.Error(err))
	}
	credential := make(chan string)
	errChan := make(chan error)

	// Create In Kong Server
	go func() {
		s.createConsumerKong(name, userId)
		result, err := s.createJWTCredential(name, key)
		if err != nil {
			global.Logger.Error("Fail to generate token", zap.Error(err))
			errChan <- err
			return
		}
		s.addGroupToConsumer(name, role)
		credential <- result
		fmt.Println("Register success")
	}()
	select {
	case cred := <-credential:
		return tokenString, cred, nil
	case err := <-errChan:
		return "", "", err
	}
}

func (s *JWTService) GenerateTokenLogin(userId string, name string, role string, credentialId string) (accessToken string, newCredentialId string, err error) {
	key, err := GenerateRandomKID()
	if err != nil {
		global.Logger.Error("Fail to generate token", zap.Error(err))
		return "", "", err
	}
	tokenString, err := s.GenerateToken(userId, name, role, key)
	if err != nil {
		global.Logger.Error("Fail to generate token", zap.Error(err))
	}
	credential := make(chan string)
	errChan := make(chan error)

	go func() {
		s.deleteJWTCredential(name, credentialId)
		result, err := s.createJWTCredential(name, key)
		if err != nil {
			global.Logger.Error("Fail to generate token", zap.Error(err))
			errChan <- err
			return
		}
		credential <- result
	}()
	select {
	case cred := <-credential:
		return tokenString, cred, nil
	case err := <-errChan:
		return "", "", err
	}
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

func (s *JWTService) createConsumerKong(username string, userId string) error {
	url := "http://kong:8001/consumers/"
	consumer := Consumer{
		Username: username,
		UserID:   userId,
	}
	jsonData, err := json.Marshal(consumer)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create consumer, status code: %d", resp.StatusCode)
	}

	fmt.Println("Consumer created successfully")
	return nil
}
func (s *JWTService) createJWTCredential(consumerID, key string) (string, error) {
	public, err := publicKeyToString(s.publicKey)
	if err != nil {
		global.Logger.Error("Fail to generate private key", zap.Error(err))
		return "", err
	}
	jwtData := map[string]string{
		"algorithm":      "RS256",
		"key":            key,
		"rsa_public_key": public,
	}
	jsonData, err := json.Marshal(jwtData)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://kong:8001/consumers/%s/jwt", consumerID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("JWT credential created successfully")

	var response CreateCredentialsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", nil
	}
	return response.Id, nil
}

func (s *JWTService) addGroupToConsumer(consumerID, groupName string) error {
	groupData := map[string]string{
		"group": groupName,
	}
	jsonData, err := json.Marshal(groupData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://kong:8001/consumers/%s/acls", consumerID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add group to consumer, status code: %d", resp.StatusCode)
	}

	fmt.Println("Group added to consumer successfully")
	return nil
}

func (s *JWTService) deleteJWTCredential(consumerID, credentialID string) error {
	url := fmt.Sprintf("http://kong:8001/consumers/%s/jwt/%s", consumerID, credentialID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete JWT credential, status code: %d", resp.StatusCode)
	}

	fmt.Println("JWT credential deleted successfully")
	return nil
}
