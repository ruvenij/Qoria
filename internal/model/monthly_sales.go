package model

type MonthlySales struct {
	Month      string `json:"month"`
	TotalSales int64  `json:"total_sales"`
}
