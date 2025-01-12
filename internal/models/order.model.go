package models

import "time"

type Order struct {
	ID         int        `db:"id" json:"id"`
	CustomerId int        `db:"customer_id" json:"customer_id"`
	Amount     float64    `db:"amount" json:"amount"`
	OrderDate  *time.Time `db:"order_date" json:"order_date"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
}

type OrderItem struct {
	ID        int     `db:"id" json:"id"`
	OrderID   int     `db:"order_id" json:"order_id"`
	ProductID int     `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
}
