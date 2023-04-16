package queries

import (
	"context"
	"fmt"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetMerchantInfo(ctx context.Context, id string) (*models.MerchantDashboardInfo, error) {
	var (
		info models.MerchantDashboardInfo
		err  error
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
	query := "UPDATE merchants SET"
	if conf.CommissionRate != nil && conf.Disabled != nil {
		query = fmt.Sprintf("%s commission=%d,disabled=%t", query,
			*conf.CommissionRate, *conf.Disabled)
	} else if conf.CommissionRate != nil {
		query = fmt.Sprintf("%s commission=%d", query, *conf.CommissionRate)
	} else if conf.Disabled != nil {
		query = fmt.Sprintf("%s disabled=%t", query, *conf.Disabled)
	}
	return db.Exec(ctx, query)
}

func GetMerchantList(ctx context.Context) ([]models.Merchant, error) {
	var merchants []models.Merchant
	rows, err := db.Query(ctx,
		"SELECT a.account_id,a.commission,a.disabled,b.username from merchants a,accounts b where a.account_id=b.account_id")
	if err != nil {
		return merchants, err
	}
	for rows.Next() {
		var temp models.Merchant
		if err := rows.Scan(&temp.AccountId, &temp.CommissionRate, &temp.Disabled,
			&temp.Name); err != nil {
			return merchants, err
		}
		merchants = append(merchants, temp)
	}
	return merchants, nil
}

func GetMerchantById(ctx context.Context, id string) (models.Merchant, error) {
	var m models.Merchant
	row, err := db.QueryRow(ctx,
		"SELECT a.account_id,a.commission,a.disabled,b.username from merchants a,accounts b "+
			"where a.account_id=b.account_id AND a.account_id=$1", id)
	if err != nil {
		return m, err
	}
	err = row.Scan(&m.AccountId, &m.CommissionRate, &m.Disabled, &m.Name)
	return m, err
}

func DeleteMerchantById(ctx context.Context, id string) error {
	return db.Exec(ctx,
		"DELETE FROM accounts WHERE account_id=$1", id)
}
