package auth

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func HashKeySecret(password string) string {
	hasher := sha3.New256()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}
