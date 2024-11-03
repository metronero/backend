package moneropay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/moneropay/go-monero/walletrpc"
	"gitlab.com/moneropay/moneropay/v2/pkg/model"

	"gitlab.com/metronero/backend/internal/utils/config"
)

func CreatePayment(amount uint64, paymentId string) (string, error) {
	cUrl, err := url.JoinPath(config.CallbackAddr, "callback", paymentId)
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

func WithdrawFunds(address string, amount uint64) (model.TransferPostResponse, error) {
	dest := walletrpc.Destination{Amount: amount, Address: address}
	tr := model.TransferPostRequest{Destinations: []walletrpc.Destination{dest}}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&tr); err != nil {
		return model.TransferPostResponse{}, err
	}

	endpoint, err := url.JoinPath(config.MoneroPay, "/transfer")
	if err != nil {
		return model.TransferPostResponse{}, err
	}
	req, err := http.NewRequest("POST", endpoint, b)
	if err != nil {
		return model.TransferPostResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: configurable timeout
	cl := &http.Client{Timeout: 15 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return model.TransferPostResponse{}, err
	}
	if resp.StatusCode != 200 {
		return model.TransferPostResponse{}, fmt.Errorf("MoneroPay /transfer failed")
	}
	var transferResp model.TransferPostResponse
	err = json.NewDecoder(resp.Body).Decode(&transferResp)
	return transferResp, err
}

// Returns health status and version of MoneroPay
func CheckHealth() (bool, string, error) {
	endpoint, err := url.JoinPath(config.MoneroPay, "/health")
	if err != nil {
		return false, "", err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return false, "", err
	}
	// TODO: configurable timeout
	cl := &http.Client{Timeout: 15 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return false, "", err
	}
	if resp.StatusCode != 200 {
		return false, "", nil
	}
	serverHeader := resp.Header.Get("Server")
	var version string
	if strings.HasPrefix(serverHeader, "MoneroPay/") {
		version = strings.TrimPrefix(serverHeader, "MoneroPay/")
	} else {
		return false, "", fmt.Errorf("Server header is missing or not MoneroPay")
	}
	return true, version, nil
}

// Returns total and unlocked wallet balance from MoneroPay
func GetBalance() (uint64, uint64, error) {
	endpoint, err := url.JoinPath(config.MoneroPay, "/balance")
	if err != nil {
		return 0, 0, err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return 0, 0, err
	}
	// TODO: configurable timeout
	cl := &http.Client{Timeout: 15 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return 0, 0, err
	}
	if resp.StatusCode != 200 {
		return 0, 0, nil
	}
	var balanceResp model.BalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&balanceResp)
	return balanceResp.Total, balanceResp.Unlocked, nil
}
