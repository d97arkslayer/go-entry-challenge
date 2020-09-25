package Repositories

import (
	"context"
	"encoding/json"
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