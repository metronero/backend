package queries

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
	db "gitlab.com/moneropay/metronero/metronero-backend/platform/database"
)

func GetPaymentsByAccount(ctx context.Context, id string) ([]models.Payment, error) {
	// If an account ID was specified, get payments that belong to this account.
	query := "SELECT id, merchant_name, amount, fee, order_id, status, last_update FROM payments"
	if id != "" {
		query = fmt.Sprintf("%s WHERE account_id=%s", query, id)
	}

	var payments []models.Payment

	rows, err := db.Query(ctx, query)
	if err != nil {
		return payments, err
	}

	for rows.Next() {
		var temp models.Payment
		if err := rows.Scan(&temp.InvoiceId, &temp.MerchantName, &temp.Amount, &temp.Fee,
		    &temp.OrderId, &temp.Status, &temp.LastUpdate); err != nil {
			// TODO: check in here whether if the error was caused unknown account_id
			// or database related error
			return payments, err
		}
		payments = append(payments, temp)
	}
	return payments, nil
}

// Get all payments from all merchants. Invoked by the admin user.
func GetPayments(ctx context.Context) ([]models.Payment, error) {
	return GetPaymentsByAccount(ctx, "")
}

func CreatePaymentRequest(ctx context.Context, merchantId, name string,
    req *models.PaymentRequest) (*models.RequestPaymentResponse, error) {
	id := uuid.New().String()
	if err := db.Exec(ctx,
	    "INSERT INTO payments(payment_id,amount,order_id,merchant_name,account_id)VALUES($1,$2,$3,$4,$5)",
	    id, req.Amount, req.OrderId, name, merchantId); err != nil {
		return nil, err
	}
	return &models.RequestPaymentResponse{PaymentId: id}, nil
}
