package models

type Payment struct {
	Id string
        Amount uint64
        AcceptUrl string
        CancelUrl string
        CallbackUrl string

	// Merchant provided ID for this payment.
        OrderId string

	// Index of subaddress that was used to accept this payment.
        AddressIndex uint64

	// Callback data from MoneroPay
        CallbackData string

	// Possible statuses: pending, confirming, finished, cancelled, withdrawn.
	Status string
}
