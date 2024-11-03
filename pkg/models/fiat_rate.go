package models

type FiatRate struct {
	Price  float64 `json:"price"`
	Source string  `json:"source"`
}
