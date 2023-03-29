package queries

import (
	"context"
	"fmt"

	"gitlab.com/metronero/backend/app/models"
	db "gitlab.com/metronero/backend/platform/database"
)

func GetWithdrawalsByAccount(ctx context.Context, id string) ([]models.Withdrawal, error) {
	// If an account ID was specified, get withdrawals that belong to this account.
	query := "SELECT withdrawal_id, amount, withdraw_date, merchant_name FROM withdrawals"
	if id != "" {
		query = fmt.Sprintf("%s WHERE account_id='%s'", query, id)
	}

	var withdrawals []models.Withdrawal

	rows, err := db.Query(ctx, query)
	if err != nil {
		return withdrawals, err
	}

	for rows.Next() {
		var temp models.Withdrawal
		if err := rows.Scan(&temp.Id, &temp.Amount, &temp.Date,
			&temp.MerchantName); err != nil {
			return withdrawals, err
		}
		withdrawals = append(withdrawals, temp)
	}
	return withdrawals, nil
}

// Get all withdrawals from all merchants. Invoked by the admin user.
func GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error) {
	return GetWithdrawalsByAccount(ctx, "")
}
