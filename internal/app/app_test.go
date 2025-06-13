package app

import (
	"Qoria/internal/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_ProcessData(t *testing.T) {
	newApp := NewApp()
	txns := []*model.Transaction{
		{
			TransactionId:   "1",
			TransactionDate: time.Now(),
			Country:         "USA",
			Region:          "California",
			ProductId:       "P001",
			ProductName:     "Product 1",
			Quantity:        2,
			TotalPrice:      decimal.NewFromFloat(120.50),
			StockQuantity:   150,
		},
		{
			TransactionId:   "2",
			TransactionDate: time.Now(),
			Country:         "India",
			Region:          "Mumbai",
			ProductId:       "P002",
			ProductName:     "Product 2",
			Quantity:        5,
			TotalPrice:      decimal.NewFromFloat(110.50),
			StockQuantity:   200,
		},
	}

	err := newApp.ProcessData(txns)
	assert.NoError(t, err)
}

func Test_GetRevenueByCountrySummary(t *testing.T) {
	newApp := NewApp()
	txns := []*model.Transaction{
		{
			TransactionId:   "1",
			TransactionDate: time.Now(),
			Country:         "USA",
			Region:          "California",
			ProductId:       "P001",
			ProductName:     "Product 1",
			Quantity:        2,
			TotalPrice:      decimal.NewFromFloat(120.50),
			StockQuantity:   150,
		},
		{
			TransactionId:   "2",
			TransactionDate: time.Now(),
			Country:         "India",
			Region:          "Mumbai",
			ProductId:       "P002",
			ProductName:     "Product 2",
			Quantity:        5,
			TotalPrice:      decimal.NewFromFloat(110.50),
			StockQuantity:   200,
		},
	}

	_ = newApp.ProcessData(txns)

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=10", nil)
	w := httptest.NewRecorder()
	newApp.GetRevenueByCountrySummary(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func Test_GetRevenueByRegionSummary(t *testing.T) {
	newApp := NewApp()
	txns := []*model.Transaction{
		{
			TransactionId:   "1",
			TransactionDate: time.Now(),
			Country:         "USA",
			Region:          "California",
			ProductId:       "P001",
			ProductName:     "Product 1",
			Quantity:        2,
			TotalPrice:      decimal.NewFromFloat(120.50),
			StockQuantity:   150,
		},
		{
			TransactionId:   "2",
			TransactionDate: time.Now(),
			Country:         "India",
			Region:          "Mumbai",
			ProductId:       "P002",
			ProductName:     "Product 2",
			Quantity:        5,
			TotalPrice:      decimal.NewFromFloat(110.50),
			StockQuantity:   200,
		},
	}

	_ = newApp.ProcessData(txns)

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=10", nil)
	w := httptest.NewRecorder()
	newApp.GetRevenueByRegionSummary(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func Test_GetProductFrequencySummary(t *testing.T) {
	newApp := NewApp()
	txns := []*model.Transaction{
		{
			TransactionId:   "1",
			TransactionDate: time.Now(),
			Country:         "USA",
			Region:          "California",
			ProductId:       "P001",
			ProductName:     "Product 1",
			Quantity:        2,
			TotalPrice:      decimal.NewFromFloat(120.50),
			StockQuantity:   150,
		},
		{
			TransactionId:   "2",
			TransactionDate: time.Now(),
			Country:         "India",
			Region:          "Mumbai",
			ProductId:       "P002",
			ProductName:     "Product 2",
			Quantity:        5,
			TotalPrice:      decimal.NewFromFloat(110.50),
			StockQuantity:   200,
		},
	}

	_ = newApp.ProcessData(txns)

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=10", nil)
	w := httptest.NewRecorder()
	newApp.GetProductFrequencySummary(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func Test_GetMonthlySalesSummary(t *testing.T) {
	newApp := NewApp()
	txns := []*model.Transaction{
		{
			TransactionId:   "1",
			TransactionDate: time.Now(),
			Country:         "USA",
			Region:          "California",
			ProductId:       "P001",
			ProductName:     "Product 1",
			Quantity:        2,
			TotalPrice:      decimal.NewFromFloat(120.50),
			StockQuantity:   150,
		},
		{
			TransactionId:   "2",
			TransactionDate: time.Now(),
			Country:         "India",
			Region:          "Mumbai",
			ProductId:       "P002",
			ProductName:     "Product 2",
			Quantity:        5,
			TotalPrice:      decimal.NewFromFloat(110.50),
			StockQuantity:   200,
		},
	}

	_ = newApp.ProcessData(txns)

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=10", nil)
	w := httptest.NewRecorder()
	newApp.GetMonthlySalesSummary(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}
