package queries

import (
	"context"
	"time"

	"github.com/google/uuid"

	db "gitlab.com/moneropay/metronero/metronero-backend/platform/database"
)

// Given an username, returns password hash
func UserLogin(ctx context.Context, username string) (string, error) {
	var passwordHash string

	row, err := db.QueryRow(ctx, "SELECT password_hash FROM accounts WHERE username=$1", username)
	if err != nil {
		return "", err
	}

	if err := row.Scan(&passwordHash); err != nil {
		return "", err
	}

	return passwordHash, nil
}

func UserRegister(ctx context.Context, username, passwordHash string) error {
	id := uuid.New().String()

	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "INSERT INTO accounts (id, username, password_hash, creation_date)" +
	    "VALUES ($1, $2, $3, $4)", id, username, passwordHash, time.Now()); err != nil {
		return err
	}
	if username != "admin" {
		if _,  err := tx.Exec(ctx, "INSERT INTO merchants (account_id) VALUES ($1)", id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "INSERT INTO merchant_stats (account_id) VALUES ($1)", id); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
