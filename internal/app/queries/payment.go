package queries

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/metronero/backend/pkg/models"
	db "gitlab.com/metronero/backend/internal/platform/database"
)

func GetPaymentsByAccount(ctx context.Context, id string) ([]models.Payment, error) {
	// If an account ID was specified, get payments that belong to this account.
	query := "SELECT payment_id,merchant_name,amount,fee,order_id,status,last_update FROM payments"
	if id != "" {
		query = fmt.Sprintf("%s WHERE account_id='%s'", query, id)
	}
	query += " ORDER BY last_update DESC"

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
func GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	return GetPaymentsByAccount(ctx, "")
}

func CreatePaymentRequest(ctx context.Context, paymentId, merchantId, name, address string,
	req *models.PostPaymentRequest) error {
	return db.Exec(ctx,
		"INSERT INTO payments(payment_id,amount,order_id,merchant_name,account_id,accept_url,"+
		"cancel_url,callback_url,merchant_extra,address,fee,last_update)"+
			"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", paymentId, req.Amount, req.OrderId,
		name, merchantId, req.AcceptUrl, req.CancelUrl, req.CallbackUrl, req.ExtraData, address, 0, time.Now())
}

func GetPaymentPageInfo(ctx context.Context, id string) (*models.PaymentPageInfo, error) {
	row, err := db.QueryRow(ctx,
		"SELECT payment_id,amount,order_id,merchant_name,accept_url,cancel_url,"+
			"address,merchant_extra,account_id,status FROM payments WHERE payment_id=$1", id)
	if err != nil {
		return nil, err
	}
	var p models.PaymentPageInfo
	if err := row.Scan(&p.InvoiceId, &p.Amount, &p.OrderId, &p.MerchantName, &p.AcceptUrl,
		&p.CancelUrl, &p.Address, &p.ExtraData, &p.TemplateId, &p.Status); err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdatePayment(ctx context.Context, id, status, data string) error {
	return db.Exec(ctx,
		"UPDATE payments SET status=$1,callback_data=$2,last_update=$3 "+
			"WHERE payment_id=$4", status, data, time.Now(), id)
}
