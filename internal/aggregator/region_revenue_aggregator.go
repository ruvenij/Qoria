package aggregator

import (
	"Qoria/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

type RegionRevenueAggregator struct {
	data map[string]map[string]*model.RegionRevenueSummary
}

func (r *RegionRevenueAggregator) Initialize() {
	r.data = make(map[string]map[string]*model.RegionRevenueSummary)
}

func (r *RegionRevenueAggregator) ProcessTransaction(tx *model.Transaction) error {
	if _, ok := r.data[tx.Region]; !ok {
		r.data[tx.Region] = make(map[string]*model.RegionRevenueSummary)
	}

	if _, ok := r.data[tx.Region][tx.ProductId]; !ok {
		r.data[tx.Region][tx.ProductId] = &model.RegionRevenueSummary{
			Region:      tx.Region,
			ProductId:   tx.ProductId,
			ProductName: tx.ProductName,
		}
	}

	// update the values
	r.data[tx.Region][tx.ProductId].TransactionCount++
	r.data[tx.Region][tx.ProductId].Revenue =
		r.data[tx.Region][tx.ProductId].Revenue.Add(tx.TotalPrice)

	return nil
}

func (r *RegionRevenueAggregator) GetResults(w http.ResponseWriter, rq *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")
	result := make([]*model.RegionRevenueSummary, 0)
	for _, products := range r.data {
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
	log.Printf("Region revenue request took %s\n", elapsed)
}
