package dto

type CompletePurchaseRequest struct {
	UserID int `json:"user_id"`
}

type CompletePurchaseResponse struct {
	OrderID    int     `json:"order_id"`
	TotalPrice float64 `json:"total_price"`
	Message    string  `json:"message"`
}
