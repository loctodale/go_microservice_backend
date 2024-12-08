package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}
func HashPasswordSalt(password string, salt string) string {
	saltedPassword := password + salt

	hashPassword := sha256.Sum256([]byte(saltedPassword))
	return hex.EncodeToString(hashPassword[:])
}

func MatchPassword(storeHash string, password string, salt string) bool {
	hashPassword := HashPasswordSalt(password, salt)
	return storeHash == hashPassword
}
