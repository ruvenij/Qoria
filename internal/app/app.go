package app

import (
	"Qoria/internal/aggregator"
	"Qoria/internal/model"
	"Qoria/internal/util"
	"log"
	"net/http"
	"sync"
)

type App struct {
	aggregators     map[int]aggregator.Aggregator
	aggregatorsLock map[int]*sync.Mutex
}

func NewApp() *App {
	aggregators := make(map[int]aggregator.Aggregator)
	aggregatorsLock := make(map[int]*sync.Mutex)
	aggregators[util.CountryRevenueAggregator] = &aggregator.CountryRevenueAggregator{}
	aggregators[util.MonthlySalesAggregator] = &aggregator.MonthlySalesAggregator{}
	aggregators[util.ProductFrequencyAggregator] = &aggregator.ProductFrequencyAggregator{}
	aggregators[util.RegionRevenueAggregator] = &aggregator.RegionRevenueAggregator{}

	for key, agg := range aggregators {
		agg.Initialize()
		aggregatorsLock[key] = &sync.Mutex{}
	}

	return &App{
		aggregators:     aggregators,
		aggregatorsLock: aggregatorsLock,
	}
}

func (app *App) ProcessData(transactions []*model.Transaction) error {
	txnChannel := make(chan *model.Transaction, 10000)
	errChannel := make(chan error)
	workerCount := 8
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for txn := range txnChannel {
				for key, agg := range app.aggregators {
					app.aggregatorsLock[key].Lock()

					err := agg.ProcessTransaction(txn)
					if err != nil {
						log.Printf("Error occurred while trying to aggregate the values for aggregator "+
							"%d, txn : %s, error : %s\n", key, txn.TransactionId, err)
						errChannel <- err
					}

					app.aggregatorsLock[key].Unlock()
				}
			}
		}()
	}

	go func() {
		for _, txn := range transactions {
			txnChannel <- txn
		}
		close(txnChannel)
	}()

	wg.Wait()
	close(errChannel)

	for err := range errChannel {
		if err != nil {
			return err
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
