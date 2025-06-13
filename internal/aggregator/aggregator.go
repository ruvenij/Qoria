package aggregator

import (
	"Qoria/internal/model"
	"net/http"
)

type Aggregator interface {
	Initialize()
	ProcessTransaction(tx *model.Transaction) error
	GetResults(w http.ResponseWriter, r *http.Request)
}
