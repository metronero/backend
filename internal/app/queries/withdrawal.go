package queries

import (
	"context"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/metronero-go/models"
)

func GetWithdrawals(ctx context.Context) ([]models.Withdrawal, error) {
	query := "SELECT withdrawal_id,amount,withdraw_date FROM withdraw ORDER BY withdraw_date DESC"

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var temp models.Withdrawal
		if err := rows.Scan(&temp.Id, &temp.Amount, &temp.Date); err != nil {
			return withdrawals, err
		}
		withdrawals = append(withdrawals, temp)
	}
	return withdrawals, nil
}

func SaveWithdrawal(ctx context.Context, w *models.Withdrawal) error {
	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx,
		"INSERT INTO withdrawals(withdrawal_id,amount,withdraw_date)"+
			"VALUES($1,$2,$3)", w.Id, w.Amount, w.Date); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, "UPDATE instance_stats SET wallet_balance=wallet_balance-$1",
		w.Amount); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
