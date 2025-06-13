package model

type ProductFrequency struct {
	ProductId              string `json:"product_id"`
	ProductName            string `json:"product_name"`
	TransactionCount       int64  `json:"transaction_count"`
	AvailableStockQuantity int64  `json:"available_stock_quantity"`
}
