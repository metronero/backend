package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
