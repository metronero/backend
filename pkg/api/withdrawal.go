package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/pkg/models"
)

func (c *ApiClient) GetMerchantWithdrawal(token string) ([]models.Withdrawal, error) {
	var p []models.Withdrawal
	resp, err := c.backendRequest(token, "GET", "/merchant/withdrawal", nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(resp, &p)
	return p, err
}

func (c *ApiClient) PostMerchantWithdrawal(token, address string) error {
	_, err := c.backendRequest(token, "POST", "/merchant/withdrawal",
		&models.WithdrawalRequest{Address: address})
	return err
}

func (c *ApiClient) GetAdminWithdrawal(token string) ([]models.Withdrawal, error) {
	var w []models.Withdrawal
	resp, err := c.backendRequest(token, "GET", "/admin/withdrawal", nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(resp, &w)
	return w, err
}
