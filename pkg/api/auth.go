package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/metronero/backend/pkg/models"
)

var (
	ErrExUnauthorized = errors.New("Username or password is wrong.")
	ErrRegisterFail = errors.New("Failed to register new user.")
)

// TODO: redo these endpoints to avoid form data, send json
func (c *ApiClient) PostLogin(username, password string) (*models.ApiTokenInfo, error) {
	endpoint, err := url.JoinPath(c.BaseUrl, "/login")
	if err != nil {
		return nil, err
	}

	resp, err := http.PostForm(endpoint, url.Values{
		"username": []string{username},
		"password": []string{password},
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrExUnauthorized
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token models.ApiTokenInfo
	err = json.Unmarshal(b, &token)
	return &token, err
}

func (c *ApiClient) PostRegister(username, password string) error {
	endpoint, err := url.JoinPath(c.BaseUrl, "/register")
	if err != nil {
		return err
	}

	resp, err := http.PostForm(endpoint, url.Values{
		"username": []string{username},
		"password": []string{password},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrRegisterFail
	}
	return nil
}
