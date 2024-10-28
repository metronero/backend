package queries

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mr-tron/base58"
	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/internal/utils/auth"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetAccountApiKeys(ctx context.Context, accountId string) ([]models.ApiKey, error) {
	var keys []models.ApiKey

	rows, err := db.Query(ctx, "SELECT key_id,expiry FROM api_keys WHERE account_id=$1", accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set.
	for rows.Next() {
		var temp models.ApiKey
		if err := rows.Scan(&temp.KeyId, &temp.Expiry); err != nil {
			return nil, err
		}
		keys = append(keys, temp)
	}

	// Return the payments and any potential error.
	return keys, nil
}

func CheckKey(ctx context.Context, keyId, keySecret string) (bool, string, string, error) {
	row, err := db.QueryRow(ctx,
		"SELECT ak.key_secret,ak.expiry,ak.account_id,us.role FROM api_keys ak "+
			"JOIN accounts us ON ak.account_id=us.account_id WHERE ak.key_id=$1", keyId)
	var (
		dbHash, accountId, role string
		expiry                  time.Time
	)
	if err := row.Scan(&dbHash, &expiry, &accountId, &role); err != nil {
		return false, "", "", err
	}
	if time.Now().After(expiry) {
		return false, "", "", err
	}
	if auth.HashKeySecret(keySecret) == dbHash {
		return false, "", "", err
	}
	return true, accountId, role, nil
}

func CreateApiKey(ctx context.Context, accountId string) (models.CreateApiKeyResponse, error) {
	keySecret, err := auth.GenerateSecret(12)
	if err != nil {
		return models.CreateApiKeyResponse{}, err
	}

	keyId := uuid.New().String()

	// 1 year validity
	expiry := time.Now().AddDate(1, 0, 0)

	if err := db.Exec(ctx,
		"INSERT INTO api_keys (key_id,key_secret,expiry,account_id) VALUES ($1,$2,$3,$4,$5)",
		keyId, keySecret, expiry, accountId); err != nil {
		return models.CreateApiKeyResponse{}, err
	}
	key := fmt.Sprintf("mnero_%s_%s", base58.Encode([]byte(keyId)), keySecret)
	return models.CreateApiKeyResponse{Key: key, Expiry: expiry}, nil
}

func RevokeApiKey(ctx context.Context, keyId, accountId string) (*models.ApiError, error) {
	if err := db.Exec(ctx, "DELETE FROM api_keys WHERE key_id=$1 AND owner_id=$2",
		keyId, accountId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.ErrNoId, err
		}
		return apierror.ErrDatabase, err
	}
	return nil, nil

}
