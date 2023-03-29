package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/metronero/backend/app/models"
)

var (
	ErrUnauthorized = errors.New("Username or password is wrong.")
	ErrRegisterFail = errors.New("Failed to register new user.")
)

func (c *ApiClient) UserLogin(username, password string) (*models.ApiTokenInfo, error) {
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
		return nil, ErrUnauthorized
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token models.ApiTokenInfo
	err = json.Unmarshal(b, &token)
	return &token, err
}

func (c *ApiClient) UserRegister(username, password string) error {
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
