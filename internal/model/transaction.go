package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	TransactionId   string          `json:"transaction_id"`
	TransactionDate time.Time       `json:"transaction_date"`
	UserId          string          `json:"user_id"`
	Country         string          `json:"country"`
	Region          string          `json:"region"`
	ProductId       string          `json:"product_id"`
	ProductName     string          `json:"product_name"`
	Category        string          `json:"category"`
	Price           decimal.Decimal `json:"price"`
	Quantity        int64           `json:"quantity"`
	TotalPrice      decimal.Decimal `json:"total_price"`
	StockQuantity   int64           `json:"stock_quantity"`
	AddedDate       time.Time       `json:"added_date"`
}
