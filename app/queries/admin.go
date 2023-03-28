package queries

import (
	"context"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetAdminInfo(ctx context.Context) (*models.AdminDashboardInfo, error) {
	var (
		a models.AdminDashboardInfo
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
