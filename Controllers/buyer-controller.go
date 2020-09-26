package Controllers

import (
	"encoding/json"
	"github.com/d97arkslayer/go-entry-challenge/Models"
	"github.com/d97arkslayer/go-entry-challenge/Repositories"
	"github.com/go-chi/chi"
	"net/http"
)

/*
 * IndexBuyers
 * Use to list all buyers
 */
func IndexBuyers(writer http.ResponseWriter, request *http.Request)  {
	buyers, err := Repositories.IndexBuyers()
	if err != nil {
		http.Error(writer,"Error obtain buyers, error: " + err.Error() ,http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(buyers)
}

/**
 * StoreBuyer
 * Use to store a new buyer
 */
func StoreBuyer(writer http.ResponseWriter, request *http.Request){
	var buyer Models.Buyer
	err := json.NewDecoder(request.Body).Decode(&buyer)
	if err != nil {
		http.Error(writer, "Error decoding body, Error: " + err.Error(), http.StatusBadRequest)
		return
	}
	var status bool
	buyer.Type = "BUYER"
	status, buyer, err = Repositories.InsertBuyer(buyer)
	if err != nil {
		http.Error(writer, "An error has ocurred when trying to insert the information in the database, Error: " + err.Error(), http.StatusBadRequest)
		return
	}
	if status == false {
		http.Error(writer, "The record could not be inserted into the database", http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(buyer)
}

/**
 * ShowBuyer
 * Use to show buyer info
 */
func ShowBuyer(writer http.ResponseWriter, request *http.Request){
	buyerId := chi.URLParam(request, "id")
	if buyerId == "" || len(buyerId) < 1 {
		http.Error(writer, "The buyer id is required", http.StatusBadRequest)
	}
	status,buyer,err := Repositories.GetBuyer(buyerId)
	if err != nil {
		http.Error(writer, "Error getting the buyer info, Error: " + err.Error(), http.StatusBadRequest)
		return
	}
	if status != true {
		http.Error(writer, "Can not get the buyer info, because de buyer does not exists", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(buyer)

}