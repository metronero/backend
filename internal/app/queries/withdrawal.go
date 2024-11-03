package queries

import (
	"context"
	"fmt"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetWithdrawals(ctx context.Context, limit int) ([]models.Withdrawal, error) {
	query := "SELECT withdrawal_id,amount,address,withdraw_date FROM withdrawals ORDER BY withdraw_date DESC"
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var temp models.Withdrawal
		if err := rows.Scan(&temp.Id, &temp.Amount, &temp.Address, &temp.Date); err != nil {
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
		"INSERT INTO withdrawals(withdrawal_id,amount,withdraw_date,address)"+
			"VALUES($1,$2,$3,$4)", w.Id, w.Amount, w.Date, w.Address); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
