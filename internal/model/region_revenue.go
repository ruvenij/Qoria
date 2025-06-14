package model

import "github.com/shopspring/decimal"

type RegionRevenueSummary struct {
	Region         string          `json:"region"`
	TotalRevenue   decimal.Decimal `json:"total_revenue"`
	TotalItemsSold int64           `json:"total_items_sold"`
}
