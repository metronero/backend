package models

import "time"

type Invoice struct {
	InvoiceId    string `json:"invoice_id"`
	MerchantName string `json:"merchant_name"`
	Amount       uint64 `json:"amount"`
	Paid         uint64 `json:"paid"`

	// Merchant provided ID for this payment
	OrderId string `json:"order_id"`

	// Possible statuses: pending, confirming, finished, cancelled, withdrawn
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"last_update"`
}

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

type PostInvoiceRequest struct {
	Amount uint64 `json:"amount"`

	// Merchant provided ID for this payment
	OrderId string `json:"order_id,omitempty"`

	// URL for redirection based on customer actions.
	// These variables are passed to the merchant template.
	AcceptUrl   string `json:"accept_url,omitempty"`
	CancelUrl   string `json:"cancel_url,omitempty"`
	CallbackUrl string `json:"callback_url,omitempty"`

	// Merchant provided extra data field
	ExtraData string `json:"extra_data,omitempty"`

	// Number of confirmations after a payment should be marked complete.
	CompleteOn *uint64 `json:"complete_on,omitempty"`

	// After how long this invoice should expire if not completed
	ExpireAfter *time.Duration `json:"expire_after,omitempty"`
}

type PostInvoiceResponse struct {
	InvoiceId string `json:"invoice_id"`
	Address   string `json:"address"`
}

type GetInvoiceCountResponse struct {
	Count   uint64 `json:"count"`
	Pending uint64 `json:"pending"`
}
