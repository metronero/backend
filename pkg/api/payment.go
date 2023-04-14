package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/pkg/models"
)

func (c *ApiClient) GetPaymentPage(paymentId string) ([]byte, error) {
	return c.backendRequest("", "GET", "/p/"+paymentId, nil)
}

func (c *ApiClient) GetMerchantPayment(token string) ([]models.Payment, error) {
	var p []models.Payment
	resp, err := c.backendRequest(token, "GET", "/merchant/payment", nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(resp, &p)
	return p, err
}

func (c *ApiClient) PostMerchantPayment(token string, req *models.PostPaymentRequest) (*models.PostPaymentResponse, error) {
	resp, err := c.backendRequest(token, "POST", "/merchant/payment", req)
	if err != nil {
		return nil, err
	}
	var p models.PostPaymentResponse
	err = json.Unmarshal(resp, &p)
	return &p, err
}

func (c *ApiClient) GetAdminPayment(token string) ([]models.Payment, error) {
	var w []models.Payment
	resp, err := c.backendRequest(token, "GET", "/admin/payment", nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(resp, &w)
	return w, err
}
