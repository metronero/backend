package queries

import (
	"context"
	"strconv"
	"strings"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func ConfigureMerchant(ctx context.Context, id string, conf *models.MerchantSettings) error {
	query := "UPDATE merchants SET "
	params := []interface{}{}
	paramIndex := 1
	if conf.CompleteOn != nil {
		query += "complete_on = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, conf.CompleteOn)
		paramIndex++
	}
	if conf.ExpireAfter != nil {
		query += "expire_after = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, conf.ExpireAfter)
		paramIndex++
	}
	if conf.FiatCurrency != nil {
		query += "fiat_currency = $" + strconv.Itoa(paramIndex) + ", "
		params = append(params, conf.FiatCurrency)
		paramIndex++
	}
	query = strings.TrimSuffix(query, ", ") + " WHERE account_id = $" + strconv.Itoa(paramIndex)
	params = append(params, id)
	return db.Exec(ctx, query, params...)
}

func GetMerchantSettings(ctx context.Context, id string) (models.MerchantSettings, error) {
	var settings models.MerchantSettings
	row, err := db.QueryRow(ctx, "SELECT complete_on,expire_after,fiat_currency FROM merchants WHERE account_id=$1", id)
	if err != nil {
		return models.MerchantSettings{}, err
	}
	if err := row.Scan(&settings.CompleteOn, &settings.ExpireAfter, &settings.FiatCurrency); err != nil {
		return models.MerchantSettings{}, err
	}
	return settings, nil
}

func GetMerchantList(ctx context.Context) ([]models.Merchant, error) {
	var merchants []models.Merchant
	rows, err := db.Query(ctx,
		"SELECT a.account_id,b.disabled,b.username from merchants a,accounts b where a.account_id=b.account_id")
	if err != nil {
		return merchants, err
	}
	for rows.Next() {
		var temp models.Merchant
		if err := rows.Scan(&temp.AccountId, &temp.Disabled,
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
		"SELECT a.account_id,a.disabled,b.username from merchants a,accounts b "+
			"where a.account_id=b.account_id AND a.account_id=$1", id)
	if err != nil {
		return m, err
	}
	err = row.Scan(&m.AccountId, &m.Disabled, &m.Name)
	return m, err
}

func DeleteMerchantById(ctx context.Context, id string) error {
	return db.Exec(ctx,
		"DELETE FROM accounts WHERE account_id=$1", id)
}

func GetMerchantCount(ctx context.Context) (uint64, uint64, error) {
	var total, active uint64
	row, err := db.QueryRow(ctx,
		`SELECT 
		COUNT(m.account_id), 
		COUNT(CASE WHEN a.disabled = false THEN 1 END) 
	FROM merchants m
	JOIN accounts a ON m.account_id = a.account_id`)
	if err != nil {
		return 0, 0, err
	}
	if err := row.Scan(&total, &active); err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func GetMerchantStats(ctx context.Context) (models.MerchantStats, error) {
	var stats models.MerchantStats
	query := `
		SELECT 
			COUNT(*) AS total_invoices,
			COUNT(CASE WHEN status = 'Pending' THEN 1 END) AS pending,
			COALESCE(SUM(amount), 0) AS total_sales
		FROM payments
	`
	row, err := db.QueryRow(ctx, query)
	if err != nil {
		return stats, err
	}
	if err := row.Scan(&stats.TotalInvoices, &stats.Pending, &stats.TotalSales); err != nil {
		return stats, err
	}
	return stats, nil
}
