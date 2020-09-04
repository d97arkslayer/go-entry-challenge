package Controllers

import (
	"encoding/json"
	"github.com/d97arkslayer/go-entry-challenge/Types"
	"net/http"
)

/*
 * IndexBuyers
 * Use to list all buyers
 */
func IndexBuyers(writer http.ResponseWriter, request *http.Request)  {
	message := Types.Message{
		Message: "Welcome to golang programming",
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(message)
}