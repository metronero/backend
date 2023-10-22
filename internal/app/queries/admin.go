package queries

import (
	"context"
	"fmt"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/metronero-go/models"
)

func GetAdminInfo(ctx context.Context) (*models.AdminDashboardInfo, error) {
	var (
		a   models.AdminDashboardInfo
		err error
	)
	if a.Stats, err = InstanceStats(ctx); err != nil {
		return nil, err
	}
	if a.RecentPayments, err = GetAllPayments(ctx); err != nil {
		return nil, err
	}
	return &a, nil
}

func AdminEditMerchant(ctx context.Context, conf *models.MerchantSettings) error {
	query := "UPDATE merchants SET"
	// TODO: support password changes here.
	if conf.Disabled != nil {
		query = fmt.Sprintf("%s disabled=%t", query, *conf.Disabled)
	}
	query += " WHERE account_id=$1"
	return db.Exec(ctx, query, conf.AccountId)
}
