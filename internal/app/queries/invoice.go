package queries

import (
	"context"
	"fmt"
	"time"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetPaymentsByAccount(ctx context.Context, id string, limit int) ([]models.Invoice, error) {
	query := `
	SELECT p.payment_id,a.username,p.amount,p.order_id,p.status,p.last_update 
	FROM payments p
	JOIN accounts a ON p.account_id=a.account_id
	`

	// If an account ID is provided, filter by account_id.
	if id != "" {
		query = fmt.Sprintf("%s WHERE p.account_id='%s'", query, id)
	}

	// Order by last_update in descending order (most recent first).
	query += " ORDER BY last_update DESC"

	// Apply the limit if it's greater than 0.
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	// Prepare to store the results.
	var payments []models.Invoice

	// Execute the query.
	rows, err := db.Query(ctx, query)
	if err != nil {
		return payments, err
	}
	defer rows.Close()

	// Iterate through the result set.
	for rows.Next() {
		var temp models.Invoice
		if err := rows.Scan(&temp.InvoiceId, &temp.MerchantName, &temp.Amount,
			&temp.OrderId, &temp.Status, &temp.LastUpdate); err != nil {
			return payments, err
		}
		payments = append(payments, temp)
	}

	// Return the payments and any potential error.
	return payments, nil
}

// Get all payments from all merchants. Invoked by the admin user.
func GetAllPayments(ctx context.Context) ([]models.Invoice, error) {
	return GetPaymentsByAccount(ctx, "", 0)
}

func CreatePaymentRequest(ctx context.Context, paymentId, merchantId, address string,
	req *models.PostInvoiceRequest) error {
	var (
		confirmations uint64
		expireAfter   time.Duration
	)
	settings, err := GetMerchantSettings(ctx, merchantId)
	if err != nil {
		return err
	}
	if req.CompleteOn == nil {
		confirmations = *settings.CompleteOn
	} else {
		confirmations = *req.CompleteOn
	}
	if req.ExpireAfter == nil {
		expireAfter = *settings.ExpireAfter
	} else {
		expireAfter = *req.ExpireAfter
	}
	now := time.Now()
	return db.Exec(ctx,
		"INSERT INTO payments(payment_id,amount,order_id,account_id,accept_url,"+
			"cancel_url,callback_url,merchant_extra,address,last_update,complete_on,expires)"+
			"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", paymentId, req.Amount,
		req.OrderId, merchantId, req.AcceptUrl, req.CancelUrl, req.CallbackUrl, req.ExtraData,
		address, now, confirmations, now.Add(expireAfter))
}

func GetPaymentPageInfo(ctx context.Context, id string) (*models.InvoicePageInfo, error) {
	row, err := db.QueryRow(ctx,
		`SELECT p.payment_id,p.amount,p.order_id,a.username,p.accept_url,p.cancel_url,
				p.address,p.merchant_extra,p.account_id,p.status,p.expires
		 FROM payments p
		 JOIN accounts a ON p.account_id=a.account_id
		 WHERE p.payment_id=$1`, id)
	if err != nil {
		return nil, err
	}
	var p models.InvoicePageInfo
	if err := row.Scan(&p.InvoiceId, &p.Amount, &p.OrderId, &p.MerchantName, &p.AcceptUrl,
		&p.CancelUrl, &p.Address, &p.ExtraData, &p.TemplateId, &p.Status, &p.Expires); err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdatePayment(ctx context.Context, id, status, data string, paid uint64) error {
	return db.Exec(ctx,
		"UPDATE invoices SET status=$1,callback_data=$2,last_update=$3,paid=$4 "+
			"WHERE payment_id=$5", status, data, time.Now(), paid, id)
}

func GetInvoiceCount(ctx context.Context) (uint64, uint64, error) {
	var total, active uint64
	row, err := db.QueryRow(ctx,
		"SELECT COUNT(*),COUNT(CASE WHEN status = 'Completed' THEN 1 END) FROM payments")
	if err != nil {
		return 0, 0, err
	}
	if err := row.Scan(&total, &active); err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func GetInvoiceCompleteOn(ctx context.Context, id string) (uint, error) {
	var completeOn uint
	row, err := db.QueryRow(ctx, "SELECT complete_on FROM payments WHERE payment_id=$1", id)
	if err != nil {
		return 0, err
	}
	if err := row.Scan(&completeOn); err != nil {
		return 0, err
	}
	return completeOn, err
}

func MarkInvoiceExpired(ctx context.Context, invoiceId string) error {
	return db.Exec(ctx, "UPDATE payments SET status='Expired' WHERE payment_id=$1", invoiceId)
}

func ExpireIncompleteInvoices() error {
	return db.Exec(context.Background(), "UPDATE payments SET status='Expired' WHERE status != 'Expired' AND status != 'Completed' AND expires < CURRENT_TIMESTAMP")
}
