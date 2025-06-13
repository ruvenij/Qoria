package main

import (
	"Qoria/api"
	"fmt"
	"net/http"
)

func main() {
	// load the data

	// process the data

	// register api functions
	api.RegisterApiFunctions()

	// start the service
	fmt.Println("Server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
