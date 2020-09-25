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
 * InsertProduct
 * This function insert a new product in DGraph
 */
func InsertProduct(product Models.Product)(bool, Models.Product, error){
	var storedProduct Models.Product
	product.Type = "PRODUCT"
	ctx := context.TODO()
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		name: string .
		price: float .
		type: string @index(exact) .
	`
	if err := dGraph.Alter(ctx, op); err != nil {
		log.Println("Error alter DGraph, Error: ", err)
		return false, storedProduct, err
	}
	mu := &api.Mutation{
		CommitNow: true,
	}

	pb, err := json.Marshal(product)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedProduct, err
	}
	mu.SetJson = pb
	response, err := dGraph.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Println("failed to marshal", err)
		return false, storedProduct, err
	}
	print("res: %v", response)
	variables := map[string]string{"$id":product.Id}
	q := `query Product($id: string){
		product(func: eq(id, $id)) {
			id
			name
			price
			type
		}
	}`
	resp, err := dGraph.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Println("Error getting the buyer, error: ", err)
		return false, storedProduct, err
	}
	type Root struct {
		Product []Models.Product `json:"product"`
	}
	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Println("Error unmarshall error: ", err)
		return false, storedProduct, err
	}
	storedProduct = r.Product[0]
	fmt.Printf("%+v\n", storedProduct)
	return true, storedProduct, nil
}

/**
 * IndexProducts
 * Use to get all products
 */
func IndexProducts() (*Types.Products, error) {
	dGraph, cancel := Database.GetDgraphClient()
	defer cancel()
	op := &api.Operation{}
	op.Schema = `
		id: string @index(exact) .
		name: string .
		price: float .
		type: string @index(exact).
		`
	ctx := context.TODO()
	errO:= dGraph.Alter(ctx, op)
	if errO != nil {
		log.Println("Error alter operation error: ", errO.Error())
		return nil, errO
	}
	q := `query products($a: string) {
		  products(func: eq(type, $a)) {
			id,
			name,
			price,
			type
		 }
		}`
	res, err := dGraph.NewTxn().QueryWithVars(ctx, q, map[string]string{"$a":"PRODUCT"})
	if err != nil {
		log.Println("Error getting the products, Error: ", err.Error())
		return nil, err
	}
	var products *Types.Products
	err = json.Unmarshal(res.Json, &products)
	if err != nil {
		log.Println("Error unmarshall the products, Error: ", err.Error())
		return nil, err
	}
	return products, nil
}