package models

import "time"

type Payment struct {
	InvoiceId    string `json:"invoice_id"`
	MerchantName string `json:"merchant_name"`
	Amount       uint64 `json:"amount"`
	AmountFloat  string `json:"amount_float,omitempty"`
	Fee          uint64 `json:"fee"`
	FeeFloat     string `json:"fee_float,omitempty"`
	// Merchant provided ID for this payment
	OrderId string `json:"order_id,omitempty"`

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

	// Merchant provided ID for this payment
	OrderId string `json:"order_id,omitempty"`

	// Possible statuses: pending, confirming, finished, cancelled, withdrawn
	Status      string    `json:"status"`
	LastUpdate  time.Time `json:"last_update"`
	AcceptUrl   string    `json:"accept_url,omitempty"`
	CancelUrl   string    `json:"cancel_url,omitempty"`
	Address     string    `json:"address"`
	Qr          string    `json:"qr"`

	// Merchant provided extra data field
	ExtraData string `json:"extra_data,omitempty"`
}

// Contains all information related to a payment
type PaymentFull struct {
	InvoiceId    string `json:"invoice_id"`
	MerchantName string `json:"merchant_name"`
	Amount       uint64 `json:"amount"`
	AmountFloat  string `json:"amount_float,omitempty"`
	Fee          uint64 `json:"fee"`
	FeeFloat     string `json:"fee_float,omitempty"`

	// Merchant provided ID for this payment
	OrderId string `json:"order_id,omitempty"`

	// Possible statuses: pending, confirming, finished, cancelled, withdrawn
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"last_update"`

	AcceptUrl   string `json:"accept_url,omitempty"`
	CancelUrl   string `json:"cancel_url,omitempty"`
	CallbackUrl string `json:"callback_url,omitempty"`
	Address     string `json:"address"`
	// Merchant provided extra data field
	ExtraData string `json:"extra_data,omitempty"`
}

type PostPaymentRequest struct {
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
}

type PostPaymentResponse struct {
	PaymentId string `json:"payment_id"`
	Address   string `json:"address"`
}
