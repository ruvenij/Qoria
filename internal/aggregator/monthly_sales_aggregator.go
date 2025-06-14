package aggregator

import (
	"Qoria/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

type MonthlySalesAggregator struct {
	data map[string]*model.MonthlySales
}

func (m *MonthlySalesAggregator) Initialize() {
	m.data = make(map[string]*model.MonthlySales)
}

// ProcessTransaction Stores the incoming value after aggregating within the data structure
func (m *MonthlySalesAggregator) ProcessTransaction(tx *model.Transaction) error {
	month := tx.TransactionDate.Month().String()
	if _, ok := m.data[month]; !ok {
		m.data[month] = &model.MonthlySales{
			Month: month,
		}
	}
	m.data[month].TotalSales += tx.Quantity
	return nil
}

// GetResults Gets the results from the data structure and returns as json
func (m *MonthlySalesAggregator) GetResults(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")
	result := make([]*model.MonthlySales, 0)
	for _, summary := range m.data {
		result = append(result, summary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalSales > result[j].TotalSales
	})

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}

	elapsed := time.Since(start)
	log.Printf("Monthly sales request took %s\n", elapsed)
}
