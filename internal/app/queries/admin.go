package queries

import (
	"context"
	"fmt"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func AdminEditMerchant(ctx context.Context, conf *models.AdminMerchantSettings, accountId string) error {
	query := "UPDATE accounts SET"
	// TODO: support password changes here.
	if conf.Disabled != nil {
		query = fmt.Sprintf("%s disabled=%t", query, *conf.Disabled)
	}
	query += " WHERE account_id=$1"
	return db.Exec(ctx, query, accountId)
}
