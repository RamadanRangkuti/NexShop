package dto

type InsertCartRequest struct {
	// User      int `json:"id" form:"id"`
	ProductID int `json:"productId" form:"product_id"`
	Quantity  int `json:"quantity" form:"quantity"`
}
