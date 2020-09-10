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
 * InsertBuyer
 * This function insert a new buyer in DGraph
 */
func InsertBuyer(buyer Models.Buyer)(bool, Models.Buyer, error){
	var storedBuyer Models.Buyer
	ctx := context.TODO()
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		name: string .
		age: int .
		type: string @index(exact) .
	`
	if err := dGraph.Alter(ctx, op); err != nil {
		log.Println("Error alter DGraph, Error: ", err)
		return false, storedBuyer, err
	}
	mu := &api.Mutation{
		CommitNow: true,
	}

	bb, err := json.Marshal(buyer)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedBuyer, err
	}
	mu.SetJson = bb
	response, err := dGraph.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedBuyer, err
	}
	print("res: %v", response)
	variables := map[string]string{"$id":buyer.Id}
	q := `query Me($id: string){
		me(func: eq(id, $id)) {
			id
			name
			age
			type
		}
	}`
	resp, err := dGraph.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Println("Error getting the buyer, error: ", err)
		return false, storedBuyer, err
	}
	type Root struct {
		Me []Models.Buyer `json:"me"`
	}
	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Println("Error unmarshall error: ", err)
		return false, storedBuyer, err
	}
	storedBuyer = r.Me[0]
	fmt.Printf("%+v\n", storedBuyer)
	return true, storedBuyer, nil
}

/**
 * IndexBuyers
 * Use to get all buyers
 */
func IndexBuyers() (*Types.Buyers, error) {
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		name: string .
		age: int .
		type: string @index(exact).
		`
	ctx := context.TODO()
	errO:= dGraph.Alter(ctx, op)
	if errO != nil {
		log.Println("Error alter operation error: ", errO.Error())
		return nil, errO
	}
	q := `query buyers($a: string) {
		  buyers(func: eq(type, $a)) {
			id,
			name,
			age,
			type
		 }
		}`
	res, err := dGraph.NewTxn().QueryWithVars(ctx, q, map[string]string{"$a":"BUYER"})
	if err != nil {
		log.Println("Error getting the buyers, Error: ", err.Error())
		return nil, err
	}
	var buyers *Types.Buyers
	err = json.Unmarshal(res.Json, &buyers)
	if err != nil {
		log.Println("Error unmarshall the buyers, Error: ", err.Error())
		return nil, err
	}
	log.Printf("%+v\n", buyers)
	return buyers, nil
}