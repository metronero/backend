package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/pkg/models"
)

func (c *ApiClient) GetAdmin(token string) (*models.AdminDashboardInfo, error) {
	resp, err := c.backendRequest(token, "GET", "/admin", nil)
	if err != nil {
		return nil, err
	}
	var a models.AdminDashboardInfo
	err = json.Unmarshal(resp, &a)
	return &a, err
}
