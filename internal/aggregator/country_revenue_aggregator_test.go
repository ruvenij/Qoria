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

func Test_CountryRevenueProcessTransaction_AddTxn(t *testing.T) {
	c := &CountryRevenueAggregator{}
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
	assert.Equal(t, 1, len(c.data["USA"]))
	assert.NotNil(t, c.data["USA"]["P001"])
	element := c.data["USA"]["P001"]
	assert.Equal(t, "USA", element.Country)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)
}

func Test_CountryRevenueProcessTransaction_TwoTxnDifferentCountriesAndProducts(t *testing.T) {
	c := &CountryRevenueAggregator{}
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
	assert.Equal(t, 1, len(c.data["USA"]))
	assert.Equal(t, 1, len(c.data["India"]))
	assert.NotNil(t, c.data["USA"]["P001"])
	assert.NotNil(t, c.data["India"]["P002"])

	element := c.data["USA"]["P001"]
	assert.Equal(t, "USA", element.Country)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)

	element = c.data["India"]["P002"]
	assert.Equal(t, "India", element.Country)
	assert.Equal(t, "P002", element.ProductId)
	assert.Equal(t, "Product 2", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(110.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)
}

func Test_CountryRevenueProcessTransaction_TwoTxnSameCountryDifferentProducts(t *testing.T) {
	c := &CountryRevenueAggregator{}
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
	assert.Equal(t, 2, len(c.data["USA"]))
	assert.NotNil(t, c.data["USA"]["P001"])
	assert.NotNil(t, c.data["USA"]["P002"])

	element := c.data["USA"]["P001"]
	assert.Equal(t, "USA", element.Country)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)

	element = c.data["USA"]["P002"]
	assert.Equal(t, "USA", element.Country)
	assert.Equal(t, "P002", element.ProductId)
	assert.Equal(t, "Product 2", element.ProductName)
	assert.Equal(t, decimal.NewFromFloat(110.50), element.Revenue)
	assert.Equal(t, int64(1), element.TransactionCount)
}

func Test_CountryRevenueProcessTransaction_TwoTxnSameCountrySameProduct(t *testing.T) {
	c := &CountryRevenueAggregator{}
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
	assert.Equal(t, 1, len(c.data["USA"]))
	assert.NotNil(t, c.data["USA"]["P001"])

	element := c.data["USA"]["P001"]
	assert.Equal(t, "USA", element.Country)
	assert.Equal(t, "P001", element.ProductId)
	assert.Equal(t, "Product 1", element.ProductName)
	assert.True(t, element.Revenue.Equal(decimal.NewFromFloat(231)))
	assert.Equal(t, int64(2), element.TransactionCount)
}

func Test_CountryRevenueGetResult_OneTxn(t *testing.T) {
	c := &CountryRevenueAggregator{}
	c.Initialize()

	_ = c.ProcessTransaction(&model.Transaction{
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

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=1", nil)
	w := httptest.NewRecorder()

	c.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.CountryRevenueSummary
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "USA", result[0].Country)
	assert.Equal(t, "P001", result[0].ProductId)
	assert.Equal(t, "Product 1", result[0].ProductName)
	assert.True(t, result[0].Revenue.Equal(decimal.NewFromFloat(120.50)))
	assert.Equal(t, int64(1), result[0].TransactionCount)
}

func Test_CountryRevenueGetResult_MoreElements(t *testing.T) {
	c := &CountryRevenueAggregator{
		data: map[string]map[string]*model.CountryRevenueSummary{
			"USA": {
				"P001": {
					Country:          "USA",
					ProductId:        "P1",
					ProductName:      "Product 1",
					TransactionCount: 2,
					Revenue:          decimal.NewFromFloat(120.50),
				},
				"P002": {
					Country:          "USA",
					ProductId:        "P2",
					ProductName:      "Product 2",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(100),
				},
			},
			"Singapore": {
				"P005": {
					Country:          "Singapore",
					ProductId:        "P5",
					ProductName:      "Product 5",
					TransactionCount: 5,
					Revenue:          decimal.NewFromFloat(1120.50),
				},
			},
			"Australia": {
				"P002": {
					Country:          "Australia",
					ProductId:        "P2",
					ProductName:      "Product 2",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(11.50),
				},
			},
			"Canada": {
				"P001": {
					Country:          "Canada",
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

	var result []*model.CountryRevenueSummary
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	assert.Equal(t, "USA", result[0].Country)
	assert.Equal(t, "P2", result[0].ProductId)
	assert.Equal(t, "Product 2", result[0].ProductName)
	assert.True(t, result[0].Revenue.Equal(decimal.NewFromFloat(100)))
	assert.Equal(t, int64(1), result[0].TransactionCount)

	assert.Equal(t, "Canada", result[1].Country)
	assert.Equal(t, "P1", result[1].ProductId)
	assert.Equal(t, "Product 1", result[1].ProductName)
	assert.True(t, result[1].Revenue.Equal(decimal.NewFromFloat(70.50)))
	assert.Equal(t, int64(1), result[1].TransactionCount)
}

func Test_CountryRevenueGetResult_MoreElementsPageExceeds(t *testing.T) {
	c := &CountryRevenueAggregator{
		data: map[string]map[string]*model.CountryRevenueSummary{
			"USA": {
				"P001": {
					Country:          "USA",
					ProductId:        "P1",
					ProductName:      "Product 1",
					TransactionCount: 2,
					Revenue:          decimal.NewFromFloat(120.50),
				},
				"P002": {
					Country:          "USA",
					ProductId:        "P2",
					ProductName:      "Product 2",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(100),
				},
			},
			"Singapore": {
				"P005": {
					Country:          "Singapore",
					ProductId:        "P5",
					ProductName:      "Product 5",
					TransactionCount: 5,
					Revenue:          decimal.NewFromFloat(1120.50),
				},
			},
			"Australia": {
				"P002": {
					Country:          "Australia",
					ProductId:        "P2",
					ProductName:      "Product 2",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(11.50),
				},
			},
			"Canada": {
				"P001": {
					Country:          "Canada",
					ProductId:        "P1",
					ProductName:      "Product 1",
					TransactionCount: 1,
					Revenue:          decimal.NewFromFloat(70.50),
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?limit=2&page=10", nil)
	w := httptest.NewRecorder()

	c.GetResults(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result []*model.CountryRevenueSummary
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))
}
