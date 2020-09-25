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
	/*err = storeBuyers(date)
	if err != nil {
		http.Error(writer, "Error getting the buyers data, Error: " + err.Error(), http.StatusBadRequest)
		return
	}*/
	err = storeProducts(date)
	if err != nil {
		http.Error(writer, "Error getting the products data, Error: " + err.Error(), http.StatusBadRequest)
	}
	writer.WriteHeader(http.StatusOK)
}


/**
 * storeProducts
 * Use this method to get products data and store in DGraph
 */
func storeProducts(date int64) error {
	productsUrl := os.Getenv("PRODUCTS_HOST") + "?date=" + strconv.FormatInt(date, 10)
	response, err := http.Get(productsUrl)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error getting buyers data, Error: "+ err.Error())
		return err
	}
	stringBody := string(body)
	lines := strings.Split(stringBody, "\n")
	for _, productString := range lines{
		var product Models.Product
		splitProduct := strings.Split(productString,"'")
		if len(splitProduct) > 1 {
			var fragmentedName [] string
			for i := 1; i <= len(splitProduct)-2; i++ {
				fragmentedName = append(fragmentedName, splitProduct[i])
			}
			name := strings.Join(fragmentedName, "'")
			name = strings.ReplaceAll(name,"\"","")
			f, err := strconv.ParseFloat(splitProduct[len(splitProduct)-1], 64)
			if err != nil {
				fmt.Println("Error parsing string to float64", err.Error())
				f = 0
			}
			product.Id = splitProduct[0]
			product.Name = name
			product.Price = f
			status, _ ,err := Repositories.InsertProduct(product)
			if err != nil {
				fmt.Println("Error inserting product, Error: "+ err.Error())
				return err
			}
			if status != true {
				fmt.Println("Can not insert product in the database")
				return errors.New("can not insert product in the database")
			}
		}
	}
	return nil
}
/**
 * storeBuyers
 * Use to get buyers data and store on DGraph database
 */
func storeBuyers(date int64) error {
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
