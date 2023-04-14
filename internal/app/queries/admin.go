package queries

import (
	"context"
	"fmt"

	"gitlab.com/metronero/backend/pkg/models"
	db "gitlab.com/metronero/backend/internal/platform/database"
)

func GetAdminInfo(ctx context.Context) (*models.AdminDashboardInfo, error) {
	var (
		a   models.AdminDashboardInfo
		err error
	)
	if a.Stats, err = InstanceStats(ctx); err != nil {
		return nil, err
	}
	if a.RecentWithdrawals, err = GetWithdrawals(ctx); err != nil {
		return nil, err
	}
	return &a, nil
}

func AdminEditMerchant(ctx context.Context, conf *models.MerchantSettings) error {
	query := "UPDATE merchants SET"
	if conf.CommissionRate != nil && conf.Disabled != nil {
		query = fmt.Sprintf("%s commission=%d,disabled=%t", query,
			*conf.CommissionRate, *conf.Disabled)
	} else if conf.CommissionRate != nil {
		query = fmt.Sprintf("%s commission=%d", query, *conf.CommissionRate)
	} else if conf.Disabled != nil {
		query = fmt.Sprintf("%s disabled=%t", query, *conf.Disabled)
	}
	query += " WHERE account_id=$1"
	return db.Exec(ctx, query, conf.AccountId)
}
