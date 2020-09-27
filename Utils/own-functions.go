package Utils

import "github.com/d97arkslayer/go-entry-challenge/Models"

/**
 * ProductContains
 * Use to validate the existence of a Product on a slice
 */
func ProductsContains(products []Models.Product, product Models.Product) bool{
	for _, item := range products {
		if item.Id == product.Id {
			return true
		}
	}
	return false
}
