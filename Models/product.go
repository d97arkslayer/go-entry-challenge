package Models

/**
 * Product
 * Model schema to product entity
 */
type Product struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Type string `json:"type"`
}