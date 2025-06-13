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

func Test_ProductFrequencyProcessTransaction_OneTxn(t *testing.T) {
	p := &ProductFrequencyAggregator{}
	p.Initialize()

	err := p.ProcessTransaction(&model.Transaction{
		TransactionId:   "1",
		TransactionDate: time.Now(),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        2,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	assert.Equal(t, 1, len(p.data))
	assert.NotNil(t, p.data["P001"])
	element := p.data["P001"]
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, int64(1), element.TransactionCount)
	assert.Equal(t, int64(150), element.AvailableStockQuantity)
}

func Test_ProductFrequencyProcessTransaction_MultipleTxnSameProduct(t *testing.T) {
	p := &ProductFrequencyAggregator{}
	p.Initialize()

	err := p.ProcessTransaction(&model.Transaction{
		TransactionId:   "1",
		TransactionDate: time.Now(),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        2,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	err = p.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Now(),
		Country:         "Germany",
		Region:          "Zelder",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        5,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	assert.Equal(t, 1, len(p.data))
	assert.NotNil(t, p.data["P001"])
	element := p.data["P001"]
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, int64(2), element.TransactionCount)
	assert.Equal(t, int64(150), element.AvailableStockQuantity)
}

func Test_ProductFrequencyProcessTransaction_MultipleTxnDifferentProducts(t *testing.T) {
	p := &ProductFrequencyAggregator{}
	p.Initialize()

	err := p.ProcessTransaction(&model.Transaction{
		TransactionId:   "1",
		TransactionDate: time.Now(),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        2,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	err = p.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Now(),
		Country:         "Germany",
		Region:          "Zelder",
		ProductId:       "P002",
		ProductName:     "Product 2",
		Quantity:        5,
		TotalPrice:      decimal.NewFromFloat(120.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	assert.Equal(t, 2, len(p.data))
	assert.NotNil(t, p.data["P001"])
	element := p.data["P001"]
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, int64(1), element.TransactionCount)
	assert.Equal(t, int64(150), element.AvailableStockQuantity)

	assert.NotNil(t, p.data["P002"])
	element = p.data["P002"]
	assert.Equal(t, "P002", element.ProductId)
	assert.Equal(t, "Product 2", element.ProductName)
	assert.Equal(t, int64(1), element.TransactionCount)
	assert.Equal(t, int64(150), element.AvailableStockQuantity)
}

func Test_ProductFrequencyGetResult_OneTxn(t *testing.T) {
	m := &ProductFrequencyAggregator{}
	m.data = map[string]*model.ProductFrequency{
		"P001": {
			ProductId:              "P001",
			ProductName:            "Product 1",
			TransactionCount:       1,
			AvailableStockQuantity: 164,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	m.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.ProductFrequency
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "P001", result[0].ProductId)
	assert.Equal(t, "Product 1", result[0].ProductName)
	assert.Equal(t, int64(1), result[0].TransactionCount)
	assert.Equal(t, int64(164), result[0].AvailableStockQuantity)
}

func Test_ProductFrequencyGetResult_MultipleTxn(t *testing.T) {
	m := &ProductFrequencyAggregator{}
	m.data = map[string]*model.ProductFrequency{
		"P001": {
			ProductId:              "P001",
			ProductName:            "Product 1",
			TransactionCount:       1,
			AvailableStockQuantity: 164,
		},
		"P005": {
			ProductId:              "P005",
			ProductName:            "Product 5",
			TransactionCount:       200,
			AvailableStockQuantity: 132,
		},
		"P004": {
			ProductId:              "P004",
			ProductName:            "Product 4",
			TransactionCount:       5,
			AvailableStockQuantity: 23,
		},
		"P006": {
			ProductId:              "P006",
			ProductName:            "Product 6",
			TransactionCount:       20,
			AvailableStockQuantity: 34,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?limit=3", nil)
	w := httptest.NewRecorder()

	m.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.ProductFrequency
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "P005", result[0].ProductId)
	assert.Equal(t, "Product 5", result[0].ProductName)
	assert.Equal(t, int64(200), result[0].TransactionCount)
	assert.Equal(t, int64(132), result[0].AvailableStockQuantity)

	assert.Equal(t, "P006", result[1].ProductId)
	assert.Equal(t, "Product 6", result[1].ProductName)
	assert.Equal(t, int64(20), result[1].TransactionCount)
	assert.Equal(t, int64(34), result[1].AvailableStockQuantity)

	assert.Equal(t, "P004", result[2].ProductId)
	assert.Equal(t, "Product 4", result[2].ProductName)
	assert.Equal(t, int64(5), result[2].TransactionCount)
	assert.Equal(t, int64(23), result[2].AvailableStockQuantity)
}
