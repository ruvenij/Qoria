package aggregator

import (
	"Qoria/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"sort"
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

	if len(result) > 30 {
		result = result[0:30]
	}

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}

	elapsed := time.Since(start)
	log.Printf("Country revenue request took %s\n", elapsed)
}
