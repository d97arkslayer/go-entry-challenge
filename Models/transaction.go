package Models

/**
 * Transaction
 * Model schema to transaction entity
 */
type Transaction struct {
	Id string `json:"id"`
	BuyerId string `json:"buyerId"`
	Ip string `json:"ip"`
	Device string `json:"device"`
	ProductIds []string `json:"productIds"`
	Type string `json:"type,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}