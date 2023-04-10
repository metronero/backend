package models

type ApiError struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Status int
}
