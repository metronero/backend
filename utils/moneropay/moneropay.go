package moneropay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/moneropay/moneropay/v2/pkg/model"

	"gitlab.com/metronero/backend/utils/config"
)

func CreatePayment(amount uint64, paymentId string) (string, error) {
	cUrl, err :=  url.JoinPath(config.CallbackAddr, "p", paymentId)
	if err != nil {
		return "", err
	}
	receiveReq := &model.ReceivePostRequest{Amount: amount, CallbackUrl: cUrl}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(receiveReq); err != nil {
		return "", err
	}

	endpoint, err := url.JoinPath(config.MoneroPay, "/receive")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", endpoint, b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: configurable timeout
	cl := &http.Client{Timeout: 15 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	var receiveResp model.ReceivePostResponse
	err = json.NewDecoder(resp.Body).Decode(&receiveResp)
	return receiveResp.Address, err
}
