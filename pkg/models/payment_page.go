package models

import "time"

// Payment data passed to checkout page template
type PaymentPageInfo struct {
	InvoiceId    string `json:"invoice_id"`
	MerchantName string `json:"merchant_name"`
	TemplateId   string `json:"template_id"`
	Amount       uint64 `json:"amount"`
	AmountFloat  string `json:"amount_float,omitempty"`
	Paid         uint64 `json:"paid"`
	PaidFloat    string `json:"paid_float,omitempty"`

	// Merchant provided ID for this payment
	OrderId string `json:"order_id,omitempty"`

	// Possible statuses: pending, confirming, finished, cancelled, withdrawn
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"last_update"`
	AcceptUrl  string    `json:"accept_url,omitempty"`
	CancelUrl  string    `json:"cancel_url,omitempty"`
	Address    string    `json:"address"`
	Qr         string    `json:"qr"`

	// Merchant provided extra data field
	ExtraData string `json:"extra_data,omitempty"`

	// When this invoice becomes invalid
	Expires time.Time `json:"expires"`
}
