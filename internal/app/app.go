package app

import (
	"Qoria/internal/aggregator"
	"Qoria/internal/model"
	"Qoria/internal/util"
	"log"
	"net/http"
)

type App struct {
	aggregators map[int]aggregator.Aggregator
}

func NewApp() *App {
	aggregators := make(map[int]aggregator.Aggregator)
	aggregators[util.CountryRevenueAggregator] = &aggregator.CountryRevenueAggregator{}
	aggregators[util.MonthlySalesAggregator] = &aggregator.MonthlySalesAggregator{}
	aggregators[util.ProductFrequencyAggregator] = &aggregator.ProductFrequencyAggregator{}
	aggregators[util.RegionRevenueAggregator] = &aggregator.RegionRevenueAggregator{}

	for _, agg := range aggregators {
		agg.Initialize()
	}

	return &App{
		aggregators: aggregators,
	}
}

func (app *App) ProcessData(transactions []*model.Transaction) error {
	for _, txn := range transactions {
		for key, agg := range app.aggregators {
			err := agg.ProcessTransaction(txn)
			if err != nil {
				log.Printf("Error occurred while trying to aggregate the values for aggregator %d, txn : %s, error : %s\n",
					key, txn.TransactionId, err)
				return err
			}
		}
	}

	return nil
}

func (app *App) GetRevenueByCountrySummary(w http.ResponseWriter, r *http.Request) {
	app.aggregators[util.CountryRevenueAggregator].GetResults(w, r)
}

func (app *App) GetRevenueByRegionSummary(w http.ResponseWriter, r *http.Request) {
	app.aggregators[util.RegionRevenueAggregator].GetResults(w, r)
}

func (app *App) GetProductFrequencySummary(w http.ResponseWriter, r *http.Request) {
	app.aggregators[util.ProductFrequencyAggregator].GetResults(w, r)
}

func (app *App) GetMonthlySalesSummary(w http.ResponseWriter, r *http.Request) {
	app.aggregators[util.MonthlySalesAggregator].GetResults(w, r)
}
