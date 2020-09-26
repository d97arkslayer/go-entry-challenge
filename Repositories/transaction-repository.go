package Repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d97arkslayer/go-entry-challenge/Database"
	"github.com/d97arkslayer/go-entry-challenge/Models"
	"github.com/d97arkslayer/go-entry-challenge/Types"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"log"
)

/**
 * InsertTransaction
 * This function insert a new transaction in DGraph
 */
func InsertTransaction(transaction Models.Transaction)(bool, Models.Transaction, error){
	var storedTransaction Models.Transaction
	transaction.Type = "TRANSACTION"
	transaction.DType = []string{"Transaction"}
	ctx := context.TODO()
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		buyerId: string @index(exact) .
		ip: string @index(exact) .
		device: string .
		productIds: [string] .
		type: string @index(exact) .

		type Transaction {
			id: string
			buyerId: string
			ip: string
			device: string
			productIds: [string]
			type: string
		}
	`
	if err := dGraph.Alter(ctx, op); err != nil {
		log.Println("Error alter DGraph, Error: ", err)
		return false, storedTransaction, err
	}
	mu := &api.Mutation{
		CommitNow: true,
	}

	tb, err := json.Marshal(transaction)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedTransaction, err
	}
	mu.SetJson = tb
	_, err = dGraph.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedTransaction, err
	}
	variables := map[string]string{"$id": transaction.Id}
	q := `query Transaction($id: string){
		transaction(func: eq(id, $id)) {
			id
			buyerId
			ip
			device
			productIds
		}
	}`
	resp, err := dGraph.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Println("Error getting the transaction, error: ", err)
		return false, storedTransaction, err
	}
	type Root struct {
		Transaction []Models.Transaction `json:"transaction"`
	}
	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Println("Error unmarshall error: ", err)
		return false, storedTransaction, err
	}
	storedTransaction = r.Transaction[0]
	return true, storedTransaction, nil
}

/**
 * IndexTransactions
 * Use to get all transactions
 */
func IndexTransactions() (*Types.Transactions, error) {
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		buyerId: string @index(exact) .
		ip: string @index(exact) .
		device: string .
		productIds: [string] .
		type: string @index(exact) .
	`
	ctx := context.TODO()
	err := dGraph.Alter(ctx, op)
	if err != nil {
		log.Println("Error alter operation error: ", err.Error())
		return nil, err
	}
	q := `query transactions($a: string) {
		  transactions(func: eq(type, $a)) {
			id,
			buyerId,
			ip,
			device,
			productIds
		 }
		}`
	res, err := dGraph.NewTxn().QueryWithVars(ctx, q, map[string]string{"$a":"TRANSACTION"})
	if err != nil {
		log.Println("Error getting the transactions, Error: ", err.Error())
		return nil, err
	}
	var transactions *Types.Transactions
	err = json.Unmarshal(res.Json, &transactions)
	if err != nil {
		log.Println("Error unmarshall the transactions, Error: ", err.Error())
		return nil, err
	}
	return transactions, nil
}

/**
 * GetTransactions
 * Use to get all buyer transactions
 */
func GetTransactions(buyerId string) ([] Models.Transaction,[] Models.Product, error){
	var transactions Types.Transactions
	var products []Models.Product
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		buyerId: string @index(exact) .
		ip: string @index(exact) .
		device: string .
		productIds: [string] .
		type: string @index(exact) .
		`
	ctx := context.TODO()
	err := dGraph.Alter(ctx, op)
	if err != nil {
		log.Println("Error alter operation error: ", err.Error())
		return transactions.Transactions,products, err
	}
	q := `query transactions($a: string) {
		  transactions(func: eq(buyerId, $a)) {
			id,
			buyerId,
			ip,
			device,
			productIds
		 }
		}`
	res, err := dGraph.NewTxn().QueryWithVars(ctx, q, map[string]string{"$a":buyerId})
	if err != nil {
		log.Println("Error getting the transactions, Error: ", err.Error())
		return transactions.Transactions, products, err
	}
	err = json.Unmarshal(res.Json, &transactions)
	if err != nil {
		log.Println("Error unmarshall the transactions, Error: ", err.Error())
		return transactions.Transactions, products, err
	}
	for _, transaction := range transactions.Transactions {
		for _, productId := range transaction.ProductIds{
			_,product, err := GetProduct(productId)
			if err != nil {
				fmt.Println("Error getting product id in shopping history, Error: " + err.Error())
				return transactions.Transactions, products, err
			}
			products = append(products, product)
		}
	}
	return transactions.Transactions, products, nil
}

/**
 * getBuyersAndProductsWithTheIP
 */
func getBuyersAndProductsWithTheIP(ip string, buyerId string)([]string, []string, error){
	var buyerIds []string
	var productIds []string
	type dataTransactions struct {
		BuyerId string `json:"buyerId"`
		ProductIds []string `json:"productIds"`
	}
	type transactionsArray struct {
		Transactions []dataTransactions `json:"transactions"`
	}
	var transactions transactionsArray
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		buyerId: string @index(exact) .
		ip: string @index(exact) .
		device: string .
		productIds: [string] .
		type: string @index(exact) .
		`
	ctx := context.TODO()
	err := dGraph.Alter(ctx, op)
	if err != nil {
		log.Println("Error alter operation error: ", err.Error())
		return buyerIds, productIds, err
	}
	q := `
	query transactions($a: string) {
  		transactionsInfo as var(func: eq(ip, "168.39.68.223")) {
    		buyerId,
    		ip,
    		productIds,
  		}

  		transactions(func: uid(transactionsInfo)) @filter(not eq(buyerId, "37a56b74"))	{
    		buyerId
    		productIds
  		}
	}
	`
	res, err := dGraph.NewTxn().QueryWithVars(ctx, q, map[string]string{"$a":ip})
	if err != nil {
		log.Println("Error getting the transactions with the same IP, Error: ", err.Error())
		return buyerIds, productIds, err
	}
	err = json.Unmarshal(res.Json, &transactions)
	if err != nil {
		log.Println("Error unmarshall the transactions with the same IP, Error: ", err.Error())
		return buyerIds, productIds, err
	}
	for _, transaction := range transactions.Transactions {
		buyerIds = append(buyerIds,transaction.BuyerId)
		productIds = append(productIds, transaction.ProductIds...)
	}
	return buyerIds, productIds, nil
}