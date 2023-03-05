package queries

import (
	"context"

	db "gitlab.com/moneropay/metronero/metronero-backend/platform/database"
)

// Given an username, returns password hash
func UserLogin(ctx context.Context, username string) (string, error) {
	var passwordHash string

	row, err := db.QueryRow(ctx, "SELECT password FROM accounts WHERE username=$1", username)
	if err != nil {
		return "", err
	}

	if err := row.Scan(&passwordHash); err != nil {
		return "", err
	}

	return passwordHash, nil
}
