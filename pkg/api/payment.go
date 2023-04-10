package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/app/models"
)

func (c *ApiClient) GetPaymentPage(paymentId string) ([]byte, error) {
	return c.backendRequest("", "GET", "/p/"+paymentId, nil)
}

func (c *ApiClient) CreatePaymentRequest(token string, req *models.PostPaymentRequest) (*models.PostPaymentResponse, error) {
	resp, err := c.backendRequest(token, "POST", "/merchant/payment", req)
	if err != nil {
		return nil, err
	}
	var p models.PostPaymentResponse
	err = json.Unmarshal(resp, &p)
	return &p, err
}
