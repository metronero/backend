package queries

import (
	"context"

	"github.com/rs/zerolog/log"

	db "gitlab.com/metronero/backend/platform/database"
)

func UpdateBalances(ctx context.Context, paymentId string, amount uint64) {
	row, err := db.QueryRow(ctx,
	    "SELECT account_id FROM payments WHERE payment_id=$1", paymentId)
	if err != nil {
		log.Error().Err(err).Str("payment_id", paymentId).Msg("Failed to get account ID from payment")
		return
	}
	var accountId string
	if err := row.Scan(&accountId); err != nil {
		log.Error().Err(err).Str("payment_id", paymentId).Msg("Failed to get account ID from payment")
		return
	}

	tx, err := db.Db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Str("payment_id", paymentId).Str("account_id", accountId).
		    Msg("Failed to update account balance")
		return
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
	    "UPDATE merchant_stats SET balance=balance+$1,total_sales=total_sales+$1 WHERE account_id=$2",
	    amount, accountId); err != nil {
		log.Error().Err(err).Str("payment_id", paymentId).Str("account_id", accountId).
		    Msg("Failed to update account balance")
		return
	}

	// TODO: change from harcoded 0 once fees are implemented
	if _, err := tx.Exec(ctx,
	    "UPDATE instance_stats SET wallet_balance=wallet_balance+$1,total_profits=$2",
	    amount, 0); err != nil {
		log.Error().Err(err).Str("payment_id", paymentId).Str("account_id", accountId).
		    Msg("Failed to update instance wallet balance")
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Str("payment_id", paymentId).Str("account_id", accountId).
		    Msg("Failed to commit balance changes")
	}
}
