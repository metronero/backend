package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/pkg/models"
)

func (c *ApiClient) GetMerchant(token string) (models.MerchantDashboardInfo, error) {
	var m models.MerchantDashboardInfo
	resp, err := c.backendRequest(token, "GET", "/merchant", nil)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(resp, &m)
	return m, err
}

func (c *ApiClient) GetAdminMerchant(token string) ([]models.Merchant, error) {
	resp, err := c.backendRequest(token, "GET", "/admin/merchant", nil)
	if err != nil {
		return nil, err
	}
	var m []models.Merchant
	err = json.Unmarshal(resp, &m)
	return m, err
}

func (c *ApiClient) GetAdminMerchantById(token, merchantId string) (models.Merchant, error) {
	var m models.Merchant
	resp, err := c.backendRequest(token, "GET", "/admin/merchant/"+merchantId, nil)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(resp, &m)
	return m, err
}

func (c *ApiClient) DeleteAdminMerchantById(token, id string) error {
	_, err := c.backendRequest(token, "DELETE", "/admin/merchant/"+id, nil)
	return err
}

func (c *ApiClient) PostAdminMerchantById(token, id string, conf interface{}) error {
	_, err := c.backendRequest(token, "POST", "/admin/merchant/"+id, conf)
	return err
}
