package Types

import "github.com/d97arkslayer/go-entry-challenge/Models"

/**
 * BuyerResponse
 * Schema to response Buyer Info
 */
type BuyerResponse struct {
	Buyer Models.Buyer `json:"buyer"`
	ShoppingHistory [] Models.Product `json:"shoppingHistory"`
	BuyersIp map[string][]Models.Buyer `json:"buyersIp"`
	ProductRecommendation []Models.Product `json:"productRecommendation"`
}
