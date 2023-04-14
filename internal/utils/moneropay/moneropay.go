package moneropay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"fmt"

	"gitlab.com/moneropay/moneropay/v2/pkg/model"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/backend/internal/utils/config"
)

func CreatePayment(amount uint64, paymentId string) (string, error) {
	cUrl, err :=  url.JoinPath(config.CallbackAddr, "callback", paymentId)
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

func WithdrawFunds(address string, amount uint64) error {
	dest := walletrpc.Destination{Amount: amount, Address: address}
	tr := model.TransferPostRequest{Destinations: []walletrpc.Destination{dest}}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&tr); err != nil {
		return err
	}

	endpoint, err := url.JoinPath(config.MoneroPay, "/transfer")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", endpoint, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: configurable timeout
	cl := &http.Client{Timeout: 15 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("MoneroPay /transfer failed")
	}
	return nil
}
