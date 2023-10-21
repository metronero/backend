package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/metronero-go/models"
)

// Given an username, returns password hash
func UserLogin(ctx context.Context, username string) (models.Account, error) {
	var loginData models.Account

	row, err := db.QueryRow(ctx, "SELECT account_id, password FROM accounts WHERE username=$1",
		username)
	if err != nil {
		return loginData, err
	}

	if err := row.Scan(&loginData.AccountId, &loginData.PasswordHash); err != nil {
		return loginData, err
	}

	go UpdateUserLastLogin(context.Background(), loginData.AccountId)

	return loginData, nil
}

func UpdateUserLastLogin(ctx context.Context, id string) {
	if err := db.Exec(ctx, "UPDATE account_stats SET last_login=$1 WHERE account_id=$2",
		time.Now(), id); err != nil {
		log.Error().Err(err).Msg("Failed to update account last login")
	}
}

func UserRegister(ctx context.Context, username, passwordHash string) error {
	id := uuid.New().String()

	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "INSERT INTO accounts (account_id, username, password)"+
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

func InvalidateToken(ctx context.Context, token string) error {
	return db.Exec(ctx, "INSERT INTO invalid_tokens (token) VALUES ($1)", token)
}
