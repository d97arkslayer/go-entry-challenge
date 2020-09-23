package Controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d97arkslayer/go-entry-challenge/Models"
	"github.com/d97arkslayer/go-entry-challenge/Repositories"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
/**

 * GetAllData
 * Use to fill the database with the data
 */
func GetAllData(writer http.ResponseWriter, request *http.Request){
	dateData := request.URL.Query().Get("date")
	if dateData == "" {
		dateData = time.Now().UTC().String()
		dateData = strings.Split(dateData, string(' '))[0]
	}
	dateTime, err := time.Parse("2006-01-02", dateData)
	if err != nil {
		http.Error(writer, "Error parsing date, Error: " + err.Error(), http.StatusBadRequest)
		return
	}
	date := dateTime.Unix()
	err = storeBuyers(date)
	if err != nil {
		http.Error(writer, "Error getting the buyers data, Error: " + err.Error(), http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

/**
 * storeBuyers
 * Use to get buyers data and store on DGraph database
 */
func storeBuyers(date int64) error{
	buyersUrl := os.Getenv("BUYERS_HOST") + "?date=" + strconv.FormatInt(date, 10)
	response, err := http.Get(buyersUrl)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error getting buyers data, Error: "+ err.Error())
		return err
	}
	var buyers [] Models.Buyer
	err = json.Unmarshal(body, &buyers)
	if err != nil {
		fmt.Println("Error unmarshal buyers data, Error: "+ err.Error())
		return err
	}
	for _, buyer := range buyers{
		status, _, err := Repositories.InsertBuyer(buyer)
		if err != nil {
			fmt.Println("Error inserting data on DGraph: "+ err.Error())
			return err
		}
		if status != true {
			fmt.Println("Can not insert data on DGraph")
			return errors.New("can not insert data on dgraph")
		}
	}
	return nil
}
