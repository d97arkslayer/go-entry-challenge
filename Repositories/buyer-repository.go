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
 * InsertBuyer
 * This function insert a new buyer in DGraph
 */
func InsertBuyer(buyer Models.Buyer)(bool, Models.Buyer, error){
	var storedBuyer Models.Buyer
	buyer.Type = "BUYER"
	buyer.DType = []string{"Buyer"}
	ctx := context.TODO()
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		name: string .
		age: int .
		type: string @index(exact) .

		type Buyer {
			id: string
			name: string
			age: int
			type: string
		}
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
	_, err = dGraph.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedBuyer, err
	}
	variables := map[string]string{"$id":buyer.Id}
	q := `query Me($id: string){
		me(func: eq(id, $id)) {
			id,
			name,
			age,
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
	err := dGraph.Alter(ctx, op)
	if err != nil {
		log.Println("Error alter operation error: ", err.Error())
		return nil, err
	}
	q := `query buyers($a: string) {
		  buyers(func: eq(type, $a)) {
			id,
			name,
			age
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
	return buyers, nil
}

/**
 * GetBuyer
 * Get buyer info
 */
func GetBuyer(id string) (bool, Models.Buyer, error) {
	var buyer Models.Buyer
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
	err := dGraph.Alter(ctx, op)
	if err != nil {
		log.Println("Error alter operation error: ", err.Error())
		return false, buyer, err
	}
	q := `
		query buyers($a: string) {
		   buyerInfo as var(func: eq(id, $a)) {
    			id,
    			name,
    			age,
    			type
			}
		buyers(func: uid(buyerInfo)) @filter(eq(type, "BUYER")) {
    			id,
    			name,
    			age
			}
		}
	`
	res, err := dGraph.NewTxn().QueryWithVars(ctx, q, map[string]string{"$a":id})
	if err != nil {
		log.Println("Error getting the buyer info, Error: ", err.Error())
		return false, buyer, err
	}
	var buyers *Types.Buyers
	err = json.Unmarshal(res.Json, &buyers)
	if len(buyers.Buyers) < 1 {
		return false, buyer, nil
	}
	if err != nil {
		log.Println("Error unmarshall the buyers, Error: ", err.Error())
		return false, buyer, err
	}
	return true, buyers.Buyers[0], nil
}