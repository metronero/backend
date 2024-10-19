package queries

import (
	"context"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func ConfigureMerchant(ctx context.Context, id string, conf *models.MerchantSettings) error {
	// TODO
	return nil
}

func GetMerchantList(ctx context.Context) ([]models.Merchant, error) {
	var merchants []models.Merchant
	rows, err := db.Query(ctx,
		"SELECT a.account_id,a.disabled,b.username from merchants a,accounts b where a.account_id=b.account_id")
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
		"SELECT COUNT(*),COUNT(CASE WHEN disabled = false THEN 1 END) FROM merchants")
	if err != nil {
		return 0, 0, err
	}
	if err := row.Scan(&total, &active); err != nil {
		return 0, 0, err
	}
	return total, active, nil
}
