package queries

import (
	"context"
	"time"

	"github.com/google/uuid"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
	db "gitlab.com/moneropay/metronero/metronero-backend/platform/database"
)

// Given an username, returns password hash
func UserLogin(ctx context.Context, username string) (models.Account, error) {
	var loginData models.Account

	row, err := db.QueryRow(ctx, "SELECT account_id, password_hash FROM accounts WHERE username=$1",
		username)
	if err != nil {
		return loginData, err
	}

	if err := row.Scan(&loginData.AccountId, &loginData.PasswordHash); err != nil {
		return loginData, err
	}

	return loginData, nil
}

func UserRegister(ctx context.Context, username, passwordHash string) error {
	id := uuid.New().String()

	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "INSERT INTO accounts (account_id, username, password_hash)"+
		"VALUES ($1, $2, $3)", id, username, passwordHash); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "INSERT INTO account_stats (account_id, creation_date)"+
		"VALUES ($1, $2)", id, time.Now()); err != nil {
		return err
	}
	if username != "admin" {
		if _, err := tx.Exec(ctx, "INSERT INTO merchants (account_id) VALUES ($1)", id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "INSERT INTO merchant_stats (account_id) VALUES ($1)", id); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
