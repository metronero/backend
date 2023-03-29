package api

import (
	"encoding/json"

	"gitlab.com/metronero/backend/app/models"
)

func (c *ApiClient) GetInstanceInfo(token string) (*models.InstanceInfo, error) {
	resp, err := c.backendRequest(token, "GET", "/admin/instance", nil)
	if err != nil {
		return nil, err
	}
	var i models.InstanceInfo
	err = json.Unmarshal(resp, &i)
	return &i, err
}

func (c *ApiClient) UpdateInstance(token string, conf *models.Instance) error {
	_, err := c.backendRequest(token, "POST", "/admin/instance", conf)
	return err
}
