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
	assert.NotNil(t, c.data["California"])
	element := c.data["California"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.TotalRevenue)
	assert.Equal(t, int64(2), element.TotalItemsSold)
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
	assert.NotNil(t, c.data["California"])
	assert.NotNil(t, c.data["Mumbai"])

	element := c.data["California"]
	assert.Equal(t, "California", element.Region)
	assert.Equal(t, decimal.NewFromFloat(120.50), element.TotalRevenue)
	assert.Equal(t, int64(2), element.TotalItemsSold)

	element = c.data["Mumbai"]
	assert.Equal(t, "Mumbai", element.Region)
	assert.Equal(t, decimal.NewFromFloat(110.50), element.TotalRevenue)
	assert.Equal(t, int64(5), element.TotalItemsSold)
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
	assert.NotNil(t, c.data["California"])

	element := c.data["California"]
	assert.Equal(t, "California", element.Region)
	assert.True(t, element.TotalRevenue.Equal(decimal.NewFromFloat(231.0)))
	assert.Equal(t, int64(7), element.TotalItemsSold)
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
	assert.NotNil(t, c.data["California"])

	element := c.data["California"]
	assert.Equal(t, "California", element.Region)
	assert.True(t, element.TotalRevenue.Equal(decimal.NewFromFloat(231)))
	assert.Equal(t, int64(7), element.TotalItemsSold)
}

func Test_RegionRevenueGetResult_OneTxn(t *testing.T) {
	c := &RegionRevenueAggregator{}
	c.data = map[string]*model.RegionRevenueSummary{
		"California": {
			Region:         "California",
			TotalRevenue:   decimal.NewFromFloat(120.50),
			TotalItemsSold: 1,
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
	assert.True(t, result[0].TotalRevenue.Equal(decimal.NewFromFloat(120.50)))
	assert.Equal(t, int64(1), result[0].TotalItemsSold)
}

func Test_RegionRevenueGetResult_MoreElements(t *testing.T) {
	c := &RegionRevenueAggregator{
		data: map[string]*model.RegionRevenueSummary{
			"California": {
				Region:         "California",
				TotalItemsSold: 2,
				TotalRevenue:   decimal.NewFromFloat(120.50),
			},
			"Geylang": {
				Region:         "Geylang",
				TotalItemsSold: 5,
				TotalRevenue:   decimal.NewFromFloat(1120.50),
			},
			"Victoria": {
				Region:         "Victoria",
				TotalItemsSold: 3,
				TotalRevenue:   decimal.NewFromFloat(11.50),
			},
			"Quebec": {
				Region:         "Quebec",
				TotalItemsSold: 1,
				TotalRevenue:   decimal.NewFromFloat(70.50),
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
	assert.True(t, result[0].TotalRevenue.Equal(decimal.NewFromFloat(1120.50)))
	assert.Equal(t, int64(5), result[0].TotalItemsSold)

	assert.Equal(t, "California", result[1].Region)
	assert.True(t, result[1].TotalRevenue.Equal(decimal.NewFromFloat(120.50)))
	assert.Equal(t, int64(2), result[1].TotalItemsSold)
}
