package api

import (
	"Qoria/internal/app"
	"net/http"
)

type Api struct {
	app *app.App
}

func NewApi(app *app.App) *Api {
	return &Api{
		app: app,
	}
}

func (a *Api) RegisterApiFunctions() {
	http.HandleFunc("/api/revenue-by-country", a.app.GetRevenueByCountrySummary)
	http.HandleFunc("/api/frequent-products", a.app.GetProductFrequencySummary)
	http.HandleFunc("/api/monthly-sales", a.app.GetMonthlySalesSummary)
	http.HandleFunc("/api/revenue-by-region", a.app.GetRevenueByRegionSummary)
}
