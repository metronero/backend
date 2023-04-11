package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/app/models"
)

func (c *ApiClient) GetMerchantList(token string) ([]models.Merchant, error) {
	resp, err := c.backendRequest(token, "GET", "/admin/merchant", nil)
	if err != nil {
		return nil, err
	}
	var m []models.Merchant
	err = json.Unmarshal(resp, &m)
	return m, err
}

func (c *ApiClient) GetMerchantById(token, merchantId string) (models.Merchant, error) {
	var m models.Merchant
	resp, err := c.backendRequest(token, "GET", "/admin/merchant/"+merchantId, nil)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(resp, &m)
	return m, err
}

func (c *ApiClient) GetMerchantInfo(token string) (models.MerchantDashboardInfo, error) {
	var m models.MerchantDashboardInfo
	resp, err := c.backendRequest(token, "GET", "/merchant", nil)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(resp, &m)
	return m, err
}

func (c *ApiClient) GetMerchantPayments(token string) ([]models.Payment, error) {
	var p []models.Payment
	resp, err := c.backendRequest(token, "GET", "/merchant/payment", nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(resp, &p)
	return p, err
}

func (c *ApiClient) GetMerchantWithdrawals(token string) ([]models.Withdrawal, error) {
	var p []models.Withdrawal
	resp, err := c.backendRequest(token, "GET", "/merchant/withdrawal", nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(resp, &p)
	return p, err
}

func (c *ApiClient) DeleteMerchantById(token, id string) error {
	_, err := c.backendRequest(token, "DELETE", "/admin/merchant/"+id, nil)
	return err
}

func (c *ApiClient) MerchantWithdrawFunds(token, address string) error {
	_, err := c.backendRequest(token, "POST", "/merchant/withdrawal",
	    &models.WithdrawalRequest{Address: address})
	return err
}
