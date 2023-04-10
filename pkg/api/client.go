package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/metronero/backend/app/models"
)

type ApiClient struct {
	Client  *http.Client
	BaseUrl string
}

func (c *ApiClient) backendRequest(token, method, endpoint string, body interface{}) ([]byte, error) {
	endpoint, err := url.JoinPath(c.BaseUrl, endpoint)
	if err != nil {
		return nil, err
	}

	b := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, endpoint, b)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// There were no API errors
	byteResp, err := io.ReadAll(res.Body)
	if res.StatusCode == 200 {
		return byteResp, err
	}

	// Parse and return the API error
	var apiErr models.ApiError
	if err = json.Unmarshal(byteResp, &apiErr); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%s (%d)", apiErr.Msg, apiErr.Code)
}
