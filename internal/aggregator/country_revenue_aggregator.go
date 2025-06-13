package aggregator

import (
	"Qoria/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type CountryRevenueAggregator struct {
	data map[string]map[string]*model.CountryRevenueSummary
}

func (c *CountryRevenueAggregator) Initialize() {
	c.data = make(map[string]map[string]*model.CountryRevenueSummary)
}

func (c *CountryRevenueAggregator) ProcessTransaction(tx *model.Transaction) error {
	if _, ok := c.data[tx.Country]; !ok {
		c.data[tx.Country] = make(map[string]*model.CountryRevenueSummary)
	}

	if _, ok := c.data[tx.Country][tx.ProductId]; !ok {
		c.data[tx.Country][tx.ProductId] = &model.CountryRevenueSummary{
			Country:     tx.Country,
			ProductId:   tx.ProductId,
			ProductName: tx.ProductName,
		}
	}

	// update the values
	c.data[tx.Country][tx.ProductId].TransactionCount++
	c.data[tx.Country][tx.ProductId].Revenue =
		c.data[tx.Country][tx.ProductId].Revenue.Add(tx.TotalPrice)

	return nil
}

func (c *CountryRevenueAggregator) GetResults(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")
	result := make([]*model.CountryRevenueSummary, 0)
	for _, products := range c.data {
		for _, summary := range products {
			result = append(result, summary)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Revenue.GreaterThanOrEqual(result[j].Revenue)
	})

	limit := 20
	page := 1

	query := r.URL.Query()
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	startVal := (page - 1) * limit
	endVal := startVal + limit

	if startVal > len(result) {
		startVal = len(result)
	}

	if endVal > len(result) {
		endVal = len(result)
	}

	result = result[startVal:endVal]
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}

	elapsed := time.Since(start)
	log.Printf("Country revenue request took %s\n", elapsed)
}
