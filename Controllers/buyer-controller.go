package Controllers

import (
	"encoding/json"
	"github.com/d97arkslayer/go-entry-challenge/Models"
	"github.com/d97arkslayer/go-entry-challenge/Repositories"
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