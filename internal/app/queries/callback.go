package queries

import (
	"context"

	"github.com/rs/zerolog/log"

	db "gitlab.com/metronero/backend/internal/platform/database"
)

func UpdateBalances(ctx context.Context, paymentId string, amount uint64) {
	row, err := db.QueryRow(ctx,
		"SELECT account_id FROM invoices WHERE invoice_id=$1", paymentId)
	if err != nil {
		log.Error().Err(err).Str("invoice_id", paymentId).Msg("Failed to get account ID from payment")
		return
	}
	var accountId string
	if err := row.Scan(&accountId); err != nil {
		log.Error().Err(err).Str("invoice_id", paymentId).Msg("Failed to get account ID from payment")
		return
	}

	tx, err := db.Db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Str("invoice_id", paymentId).Str("account_id", accountId).
			Msg("Failed to update account balance")
		return
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		"UPDATE merchants SET total_sales=total_sales+$1 WHERE account_id=$2",
		amount, accountId); err != nil {
		log.Error().Err(err).Str("invoice_id", paymentId).Str("account_id", accountId).
			Msg("Failed to update merchant's total sales")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Str("invoice_id", paymentId).Str("account_id", accountId).
			Msg("Failed to commit balance changes")
	}
}
