package queries

import (
	"context"
	"fmt"

	"gitlab.com/metronero/backend/app/models"
	db "gitlab.com/metronero/backend/platform/database"
)

func GetWithdrawalsByAccount(ctx context.Context, id string) ([]models.Withdrawal, error) {
	// If an account ID was specified, get withdrawals that belong to this account.
	query := "SELECT withdrawal_id,amount,withdraw_date,merchant_name FROM withdrawals"
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

func SaveWithdrawal(ctx context.Context, w *models.Withdrawal) error {
	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx,
	    "INSERT INTO withdrawals(withdrawal_id,amount,withdraw_date,account_id,merchant_name)"+
	    "VALUES($1,$2,$3,$4,$5)", w.Id, w.Amount, w.Date, w.AccountId, w.MerchantName);
	    err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
	    "UPDATE merchant_stats SET balance=balance-$1 WHERE account_id=$2",
	    w.Amount, w.AccountId); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, "UPDATE instance_stats SET wallet_balance=wallet_balance-$1",
	    w.Amount); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func GetWithdrawableAmount(ctx context.Context, accountId string) (uint64, error) {
	row, err := db.QueryRow(ctx,
		"SELECT balance FROM merchant_stats WHERE account_id=$1", accountId)
	if err != nil {
		return 0, err
	}
	var amount uint64
	err = row.Scan(&amount)
	return amount, err
}
