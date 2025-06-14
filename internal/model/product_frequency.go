package model

import "time"

type ProductFrequency struct {
	ProductId              string    `json:"product_id"`
	ProductName            string    `json:"product_name"`
	TransactionCount       int64     `json:"transaction_count"`
	AvailableStockQuantity int64     `json:"available_stock_quantity"`
	UnitsSold              int64     `json:"units_sold"`
	StockAddedDate         time.Time `json:"stock_added_date"`
}
