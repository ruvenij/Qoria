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

type ProductFrequencyAggregator struct {
	data map[string]*model.ProductFrequency
}

func (p *ProductFrequencyAggregator) Initialize() {
	p.data = make(map[string]*model.ProductFrequency)
}

func (p *ProductFrequencyAggregator) ProcessTransaction(tx *model.Transaction) error {
	if _, ok := p.data[tx.ProductId]; !ok {
		p.data[tx.ProductId] = &model.ProductFrequency{
			ProductId:              tx.ProductId,
			ProductName:            tx.ProductName,
			AvailableStockQuantity: tx.StockQuantity,
		}
	}

	p.data[tx.ProductId].TransactionCount++
	p.data[tx.ProductId].UnitsSold += tx.Quantity
	return nil
}

func (p *ProductFrequencyAggregator) GetResults(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")
	result := make([]*model.ProductFrequency, 0)
	for _, summary := range p.data {
		result = append(result, summary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UnitsSold > result[j].UnitsSold
	})

	limit := 20
	query := r.URL.Query()
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
	log.Printf("Product frequency request took %s\n", elapsed)
}
