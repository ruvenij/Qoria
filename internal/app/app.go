package app

import (
	"Qoria/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

type App struct {
	countryRevenueAggregate   map[string]map[string]*model.CountryRevenueSummary
	regionRevenueAggregate    map[string]map[string]*model.RegionRevenueSummary
	productFrequencyAggregate map[string]*model.ProductFrequency
	monthlySalesAggregate     map[string]*model.MonthlySales
}

func NewApp() *App {
	return &App{
		countryRevenueAggregate:   make(map[string]map[string]*model.CountryRevenueSummary),
		regionRevenueAggregate:    make(map[string]map[string]*model.RegionRevenueSummary),
		productFrequencyAggregate: make(map[string]*model.ProductFrequency),
		monthlySalesAggregate:     make(map[string]*model.MonthlySales),
	}
}

func (app *App) ProcessData(transactions []*model.Transaction) error {
	for _, transaction := range transactions {
		// process for country revenue
		err := app.processCountryRevenue(transaction)
		if err != nil {
			log.Println("Error when processing country revenue, error : ", err)
			return err
		}

		// process for region revenue
		err = app.processRegionRevenue(transaction)
		if err != nil {
			log.Println("Error when processing region revenue, error : ", err)
			return err
		}

		// process for product frequency
		err = app.processProductFrequency(transaction)
		if err != nil {
			log.Println("Error when processing product frequency, error : ", err)
			return err
		}

		// process for monthly sales
		err = app.processMonthlySales(transaction)
		if err != nil {
			log.Println("Error when processing monthly sales, error : ", err)
			return err
		}
	}

	return nil
}

func (app *App) processCountryRevenue(tx *model.Transaction) error {
	if _, ok := app.countryRevenueAggregate[tx.Country]; !ok {
		app.countryRevenueAggregate[tx.Country] = make(map[string]*model.CountryRevenueSummary)
		app.countryRevenueAggregate[tx.Country][tx.ProductId] = &model.CountryRevenueSummary{
			Country:     tx.Country,
			ProductId:   tx.ProductId,
			ProductName: tx.ProductName,
		}
	}

	if _, ok := app.countryRevenueAggregate[tx.Country][tx.ProductId]; !ok {
		app.countryRevenueAggregate[tx.Country][tx.ProductId] = &model.CountryRevenueSummary{
			Country:     tx.Country,
			ProductId:   tx.ProductId,
			ProductName: tx.ProductName,
		}
	}

	// update the values
	app.countryRevenueAggregate[tx.Country][tx.ProductId].TransactionCount++
	app.countryRevenueAggregate[tx.Country][tx.ProductId].Revenue =
		app.countryRevenueAggregate[tx.Country][tx.ProductId].Revenue.Add(tx.TotalPrice)

	return nil
}

func (app *App) processRegionRevenue(tx *model.Transaction) error {
	if _, ok := app.regionRevenueAggregate[tx.Region]; !ok {
		app.regionRevenueAggregate[tx.Region] = make(map[string]*model.RegionRevenueSummary)
		app.regionRevenueAggregate[tx.Region][tx.ProductId] = &model.RegionRevenueSummary{
			Region:      tx.Region,
			ProductId:   tx.ProductId,
			ProductName: tx.ProductName,
		}
	}

	if _, ok := app.regionRevenueAggregate[tx.Region][tx.ProductId]; !ok {
		app.regionRevenueAggregate[tx.Region][tx.ProductId] = &model.RegionRevenueSummary{
			Region:      tx.Region,
			ProductId:   tx.ProductId,
			ProductName: tx.ProductName,
		}
	}

	// update the values
	app.regionRevenueAggregate[tx.Region][tx.ProductId].TransactionCount++
	app.regionRevenueAggregate[tx.Region][tx.ProductId].Revenue =
		app.regionRevenueAggregate[tx.Region][tx.ProductId].Revenue.Add(tx.TotalPrice)

	return nil
}

func (app *App) processProductFrequency(tx *model.Transaction) error {
	if _, ok := app.productFrequencyAggregate[tx.ProductId]; !ok {
		app.productFrequencyAggregate[tx.ProductId] = &model.ProductFrequency{
			ProductId:              tx.ProductId,
			ProductName:            tx.ProductName,
			AvailableStockQuantity: tx.StockQuantity,
		}
	}

	app.productFrequencyAggregate[tx.ProductId].TransactionCount++
	return nil
}

func (app *App) processMonthlySales(tx *model.Transaction) error {
	month := tx.TransactionDate.Month().String()
	if _, ok := app.monthlySalesAggregate[month]; !ok {
		app.monthlySalesAggregate[month] = &model.MonthlySales{
			Month: month,
		}
	}
	app.monthlySalesAggregate[month].TotalSales += tx.Quantity
	return nil
}

func (app *App) GetRevenueByCountrySummary(w http.ResponseWriter, r *http.Request) {
	result := make([]*model.CountryRevenueSummary, 0)
	for _, products := range app.countryRevenueAggregate {
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
}

func (app *App) GetRevenueByRegionSummary(w http.ResponseWriter, r *http.Request) {
	result := make([]*model.RegionRevenueSummary, 0)
	for _, products := range app.regionRevenueAggregate {
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
}

func (app *App) GetProductFrequencySummary(w http.ResponseWriter, r *http.Request) {
	result := make([]*model.ProductFrequency, 0)
	for _, summary := range app.productFrequencyAggregate {
		result = append(result, summary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TransactionCount > result[j].TransactionCount
	})

	if len(result) > 20 {
		result = result[0:20]
	}

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}
}

func (app *App) GetMonthlySalesSummary(w http.ResponseWriter, r *http.Request) {
	result := make([]*model.MonthlySales, 0)
	for _, summary := range app.monthlySalesAggregate {
		result = append(result, summary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalSales > result[j].TotalSales
	})

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}
}
