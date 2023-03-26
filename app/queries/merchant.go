package queries

import (
	"context"

	db "gitlab.com/moneropay/metronero/metronero-backend/platform/database"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetMerchantInfo(ctx context.Context, id string) (*models.MerchantDashboardInfo, error) {
	var (
		info models.MerchantDashboardInfo
		err error
	)
	info.Stats, err = MerchantStats(ctx, id)
	if err != nil {
		return nil, err
	}
	info.RecentPayments, err = GetPaymentsByAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func MerchantStats(ctx context.Context, id string) (models.MerchantStats, error) {
	var stats models.MerchantStats
	row, err := db.QueryRow(ctx,
	    "SELECT balance, total_sales FROM merchant_stats WHERE account_id=$1", id)
	if err != nil {
		return stats, err
	}
	if err := row.Scan(&stats.Balance, &stats.TotalSales); err != nil {
		return stats, err
	}
	return stats, nil
}

func ConfigureMerchant(ctx context.Context, id string, conf *models.MerchantSettings) error {
	return db.Exec(ctx,
	    "UPDATE merchants SET commission = $1 AND disabled = $2 WHERE account_id = $3",
	    conf.CommissionRate, conf.Disabled, id)
}
