package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// For generating temporary user passwords. It is recommended to change
// this password after logging in.
func GeneratePassword() (string, error) {
	// 10 character long password
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
