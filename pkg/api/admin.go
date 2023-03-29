package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/app/models"
)

func (c *ApiClient) GetAdminDashboard(token string) (*models.AdminDashboardInfo, error) {
	resp, err := c.backendRequest(token, "GET", "/admin", nil)
	if err != nil {
		return nil, err
	}
	var a models.AdminDashboardInfo
	err = json.Unmarshal(resp, &a)
	return &a, err
}

func (c *ApiClient) AdminGetWithdrawals(token string) ([]models.Withdrawal, error) {
	var w []models.Withdrawal
	resp, err := c.backendRequest(token, "GET", "/admin/withdrawal", nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(resp, &w)
	return w, err
}

func (c *ApiClient) AdminGetPayments(token string) ([]models.Payment, error) {
	var w []models.Payment
	resp, err := c.backendRequest(token, "GET", "/admin/payment", nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(resp, &w)
	return w, err
}

func (c *ApiClient) AdminEditMerchant(token, id string, conf interface{}) error {
	_, err := c.backendRequest(token, "POST", "/admin/merchant/"+id, conf)
	return err
}
