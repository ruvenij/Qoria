package model

import "github.com/shopspring/decimal"

type CountryRevenueSummary struct {
	Country          string          `json:"country"`
	ProductId        string          `json:"product_id"`
	ProductName      string          `json:"product_name"`
	Revenue          decimal.Decimal `json:"revenue"`
	TransactionCount int64           `json:"transaction_count"`
}
