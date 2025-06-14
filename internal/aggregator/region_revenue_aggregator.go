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

type RegionRevenueAggregator struct {
	data map[string]*model.RegionRevenueSummary
}

func (r *RegionRevenueAggregator) Initialize() {
	r.data = make(map[string]*model.RegionRevenueSummary)
}

// ProcessTransaction Stores the incoming value after aggregating within the data structure
func (r *RegionRevenueAggregator) ProcessTransaction(tx *model.Transaction) error {
	if _, ok := r.data[tx.Region]; !ok {
		r.data[tx.Region] = &model.RegionRevenueSummary{
			Region: tx.Region,
		}
	}

	// update the values
	r.data[tx.Region].TotalItemsSold += tx.Quantity
	r.data[tx.Region].TotalRevenue = r.data[tx.Region].TotalRevenue.Add(tx.TotalPrice)

	return nil
}

// GetResults Gets the results from the data structure and returns as json
func (r *RegionRevenueAggregator) GetResults(w http.ResponseWriter, rq *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")
	result := make([]*model.RegionRevenueSummary, 0)
	for _, summary := range r.data {
		result = append(result, summary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalRevenue.GreaterThanOrEqual(result[j].TotalRevenue)
	})

	limit := 30
	query := rq.URL.Query()
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	startVal := 0
	endVal := startVal + limit
	if endVal > len(result) {
		endVal = len(result)
	}

	result = result[startVal:endVal]

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}

	elapsed := time.Since(start)
	log.Printf("Region revenue request took %s\n", elapsed)
}
