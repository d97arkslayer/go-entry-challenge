package Controllers

import (
	"encoding/json"
	"github.com/d97arkslayer/go-entry-challenge/Repositories"
	"net/http"
)

/*
 * IndexBuyers
 * Use to list all buyers
 */
func IndexProducts(writer http.ResponseWriter, request *http.Request)  {
	products, err := Repositories.IndexProducts()
	if err != nil {
		http.Error(writer,"Error obtain products, error: " + err.Error() ,http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(products)
}
