package auth

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/metronero/backend/internal/utils/config"
)

func CreateUserToken(username, id string, lifetime time.Duration) (string, time.Time, error) {
	claims := map[string]interface{}{"username": username, "id": id}
	expiryDate := time.Unix(jwtauth.ExpireIn(lifetime), 0)
	jwtauth.SetExpiry(claims, expiryDate)
	_, token, err := config.JwtSecret.Encode(claims)
	if err != nil {
		return "", expiryDate, err
	}
	return token, expiryDate, nil
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
