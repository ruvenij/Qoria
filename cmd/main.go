package main

import (
	"Qoria/internal/api"
	"Qoria/internal/app"
	"Qoria/internal/data"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	start := time.Now()
	// load the data
	transactions, err := data.LoadCsvData("data/GO_test_5m.csv")
	if err != nil {
		log.Fatal("Failed to load csv data ", err)
		return
	}

	// process the data
	newApp := app.NewApp()
	err = newApp.ProcessData(transactions)
	if err != nil {
		log.Fatal("Failed to process data ", err)
		return
	}

	seconds := time.Since(start).Seconds()
	fmt.Printf("Load time : %.2f seconds\n", seconds)

	// register api functions
	newApi := api.NewApi(newApp)
	newApi.RegisterApiFunctions()

	// start the service
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Server starting at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
