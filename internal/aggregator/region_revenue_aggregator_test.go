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

func Test_RegionRevenueProcessTransaction_AddTxn(t *testing.T) {
	c := &RegionRevenueAggregator{}
	c.Initialize()

	err := c.ProcessTransaction(&model.Transaction{
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

	// validate the entry added
	assert.Equal(t, 1, len(c.data))
	assert.Equal(t, 1, len(c.data["California"]))
	assert.NotNil(t, c.data["California"]["P001"])
	element := c.data["California"]["P001"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)
}

func Test_RegionRevenueProcessTransaction_TwoTxnDifferentCountriesAndProducts(t *testing.T) {
	c := &RegionRevenueAggregator{}
	c.Initialize()

	err := c.ProcessTransaction(&model.Transaction{
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

	err = c.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Now(),
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
	assert.Equal(t, 2, len(c.data))
	assert.Equal(t, 1, len(c.data["California"]))
	assert.Equal(t, 1, len(c.data["Mumbai"]))
	assert.NotNil(t, c.data["California"]["P001"])
	assert.NotNil(t, c.data["Mumbai"]["P002"])

	element := c.data["California"]["P001"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)

	element = c.data["Mumbai"]["P002"]
	assert.Equal(t, "Mumbai", element.Region)
	assert.Equal(t, "P002", element.ProductId)
	assert.Equal(t, "Product 2", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(110.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)
}

func Test_RegionRevenueProcessTransaction_TwoTxnSameRegionDifferentProducts(t *testing.T) {
	c := &RegionRevenueAggregator{}
	c.Initialize()

	err := c.ProcessTransaction(&model.Transaction{
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

	err = c.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Now(),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P002",
		ProductName:     "Product 2",
		Quantity:        5,
		TotalPrice:      decimal.NewFromFloat(110.50),
		StockQuantity:   200,
	})

	assert.NoError(t, err)

	// validate the entry added
	assert.Equal(t, 1, len(c.data))
	assert.Equal(t, 2, len(c.data["California"]))
	assert.NotNil(t, c.data["California"]["P001"])
	assert.NotNil(t, c.data["California"]["P002"])

	element := c.data["California"]["P001"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)

	element = c.data["California"]["P002"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, "P002", element.ProductId)
	assert.Equal(t, "Product 2", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(110.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)
}

func Test_RegionRevenueProcessTransaction_TwoTxnSameRegionSameProduct(t *testing.T) {
	c := &RegionRevenueAggregator{}
	c.Initialize()

	err := c.ProcessTransaction(&model.Transaction{
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

	err = c.ProcessTransaction(&model.Transaction{
		TransactionId:   "2",
		TransactionDate: time.Now(),
		Country:         "USA",
		Region:          "California",
		ProductId:       "P001",
		ProductName:     "Product 1",
		Quantity:        5,
		TotalPrice:      decimal.NewFromFloat(110.50),
		StockQuantity:   150,
	})

	assert.NoError(t, err)

	// validate the entry added
	assert.Equal(t, 1, len(c.data))
	assert.Equal(t, 1, len(c.data["California"]))
	assert.NotNil(t, c.data["California"]["P001"])

	element := c.data["California"]["P001"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.True(t, element.Revenue.Equal(decimal.NewFromFloat(231)))
	assert.Equal(t, int64(2), element.TransactionCount)
}

func Test_RegionRevenueGetResult_OneTxn(t *testing.T) {
	c := &RegionRevenueAggregator{}
	c.data = map[string]map[string]*model.RegionRevenueSummary{
		"California": {
			"P001": {
				Region:           "California",
				ProductId:        "P001",
				ProductName:      "Product 1",
				Revenue:          decimal.NewFromFloat(120.50),
				TransactionCount: 1,
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?limit=2", nil)
	w := httptest.NewRecorder()

	c.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.RegionRevenueSummary
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "California", result[0].Region)
	assert.Equal(t, "P001", result[0].ProductId)
	assert.Equal(t, "Product 1", result[0].ProductName)
	assert.True(t, result[0].Revenue.Equal(decimal.NewFromFloat(120.50)))
	assert.Equal(t, int64(1), result[0].TransactionCount)
}

func Test_RegionRevenueGetResult_MoreElements(t *testing.T) {
	c := &RegionRevenueAggregator{
		data: map[string]map[string]*model.RegionRevenueSummary{
			"California": {
				"P001": {
					Region:           "California",
					ProductId:        "P1",
					ProductName:      "Product 1",
					TransactionCount: 2,
					Revenue:          decimal.NewFromFloat(120.50),
				},
				"P002": {
					Region:           "California",
					ProductId:        "P2",
					ProductName:      "Product 2",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(100),
				},
			},
			"Geylang": {
				"P005": {
					Region:           "Geylang",
					ProductId:        "P5",
					ProductName:      "Product 5",
					TransactionCount: 5,
					Revenue:          decimal.NewFromFloat(1120.50),
				},
			},
			"Victoria": {
				"P002": {
					Region:           "Victoria",
					ProductId:        "P2",
					ProductName:      "Product 2",
					TransactionCount: 3,
					Revenue:          decimal.NewFromFloat(11.50),
				},
			},
			"Quebec": {
				"P001": {
					Region:           "Quebec",
					ProductId:        "P1",
					ProductName:      "Product 1",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(70.50),
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=2", nil)
	w := httptest.NewRecorder()

	c.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.RegionRevenueSummary
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	assert.Equal(t, "Geylang", result[0].Region)
	assert.Equal(t, "P5", result[0].ProductId)
	assert.Equal(t, "Product 5", result[0].ProductName)
	assert.True(t, result[0].Revenue.Equal(decimal.NewFromFloat(1120.50)))
	assert.Equal(t, int64(5), result[0].TransactionCount)

	assert.Equal(t, "California", result[1].Region)
	assert.Equal(t, "P1", result[1].ProductId)
	assert.Equal(t, "Product 1", result[1].ProductName)
	assert.True(t, result[1].Revenue.Equal(decimal.NewFromFloat(120.50)))
	assert.Equal(t, int64(2), result[1].TransactionCount)
}
