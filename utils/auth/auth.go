package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/go-chi/jwtauth/v5"

	"gitlab.com/moneropay/metronero/metronero-backend/utils/config"
)

func CreateUserToken(username string, lifetime time.Duration) (string, string, error) {
	claims := map[string]interface{}{"username": username}
	expiryDate := time.Unix(jwtauth.ExpireIn(lifetime), 0)
	jwtauth.SetExpiry(claims, expiryDate)
	_, token, err := config.JwtSecret.Encode(claims)
	if err != nil {
		return "", "", err
	}
	return token, expiryDate.Format(time.RFC3339), nil
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
