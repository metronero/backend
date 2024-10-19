package models

type GetInstanceInfoResponse struct {
	Version          string `json:"version"`
	MoneroPayVersion string `json:"moneropay"`
	Healthy          bool   `json:"healthy"`
}
