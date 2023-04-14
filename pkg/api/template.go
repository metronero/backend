package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (c *ApiClient) GetMerchantTemplate(token string) ([]byte, error) {
	return c.backendRequest(token, "GET", "/merchant/template", nil)
}

func (c *ApiClient) PostMerchantTemplate(token string, template io.Reader) error {
	var (
		b   bytes.Buffer
		err error
		fw  io.Writer
	)
	w := multipart.NewWriter(&b)
	if fw, err = w.CreateFormFile("file", "template"); err != nil {
		return err
	}
	if _, err = io.Copy(fw, template); err != nil {
		return err
	}
	w.Close()

	endpoint, err := url.JoinPath(c.BaseUrl, "/merchant/template")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	_, err = c.Client.Do(req)
	return err
}

func (c *ApiClient) DeleteMerchantTemplate(token string) error {
	_, err := c.backendRequest(token, "DELETE", "/merchant/template", nil)
	return err
}
