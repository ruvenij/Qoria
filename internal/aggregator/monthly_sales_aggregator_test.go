package aggregator

import (
	"Qoria/internal/model"
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_MonthlyRevenueProcessTransaction_AddTxn(t *testing.T) {
	m := &MonthlySalesAggregator{}
	m.Initialize()

	err := m.ProcessTransaction(&model.Transaction{
		TransactionId:   "1",
		TransactionDate: time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        2,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	assert.Equal(t, 1, len(m.data))
	assert.NotNil(t, m.data["March"])
	element := m.data["March"]
	assert.Equal(t, "March", element.Month)
	assert.Equal(t, int64(2), element.TotalSales)
}

func Test_MonthlyRevenueProcessTransaction_TwoTxnDifferentMonths(t *testing.T) {
	m := &MonthlySalesAggregator{}
	m.Initialize()

	err := m.ProcessTransaction(&model.Transaction{
		TransactionId:   "1",
		TransactionDate: time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        2,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	err = m.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Date(2021, time.November, 1, 0, 0, 0, 0, time.UTC),
		Country:         "India",
		Region:          "Mumbai",
		ProductId:       "P002",
		ProductName:     "Product 2",
		Quantity:        5,
		TotalPrice:      decimal.NewFromFloat(110.50),
		StockQuantity:   200,
	})

	assert.NoError(t, err)

	// validate the entry added
	assert.Equal(t, 2, len(m.data))
	assert.NotNil(t, m.data["March"])
	assert.NotNil(t, m.data["November"])

	element := m.data["March"]
	assert.Equal(t, "March", element.Month)
	assert.Equal(t, int64(2), element.TotalSales)

	element = m.data["November"]
	assert.Equal(t, "November", element.Month)
	assert.Equal(t, int64(5), element.TotalSales)
}

func Test_MonthlyRevenueProcessTransaction_TwoTxnDSameMonth(t *testing.T) {
	m := &MonthlySalesAggregator{}
	m.Initialize()

	err := m.ProcessTransaction(&model.Transaction{
		TransactionId:   "1",
		TransactionDate: time.Date(2021, time.November, 25, 0, 0, 0, 0, time.UTC),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        2,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	err = m.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
		Country:         "India",
		Region:          "Mumbai",
		ProductId:       "P002",
		ProductName:     "Product 2",
		Quantity:        5,
		TotalPrice:      decimal.NewFromFloat(110.50),
		StockQuantity:   200,
	})

	assert.NoError(t, err)

	// validate the entry added
	assert.Equal(t, 1, len(m.data))
	assert.NotNil(t, m.data["November"])

	element := m.data["November"]
	assert.Equal(t, "November", element.Month)
	assert.Equal(t, int64(7), element.TotalSales)
}

func Test_MonthlyRevenueGetResult_OneTxn(t *testing.T) {
	m := &MonthlySalesAggregator{}
	m.data = map[string]*model.MonthlySales{
		"March": {
			TotalSales: 120,
			Month:      "March",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=1", nil)
	w := httptest.NewRecorder()

	m.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.MonthlySales
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "March", result[0].Month)
	assert.Equal(t, int64(120), result[0].TotalSales)
}

func Test_MonthlyRevenueGetResult_MultipleResults(t *testing.T) {
	m := &MonthlySalesAggregator{}
	m.data = map[string]*model.MonthlySales{
		"March": {
			Month:      "March",
			TotalSales: 120,
		},
		"May": {
			Month:      "May",
			TotalSales: 210,
		},
		"November": {
			Month:      "November",
			TotalSales: 50,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=1", nil)
	w := httptest.NewRecorder()

	m.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.MonthlySales
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "May", result[0].Month)
	assert.Equal(t, int64(210), result[0].TotalSales)

	assert.Equal(t, "March", result[1].Month)
	assert.Equal(t, int64(120), result[1].TotalSales)

	assert.Equal(t, "November", result[2].Month)
	assert.Equal(t, int64(50), result[2].TotalSales)
}
