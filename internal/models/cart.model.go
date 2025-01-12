package models

import "time"

type ShoppingCart struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	ProductID int       `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DetailShoppingCart struct {
	ID          int    `db:"id" json:"id"`
	UserName    string `db:"username" json:"-"`
	ProductID   int    `db:"product_id" json:"product_id"`
	ProductName string `db:"name" json:"productName"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

type DetailItemShoppingCart struct {
	ID          int    `db:"id" json:"id"`
	ProductName string `db:"name" json:"name"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

type InsertShopping struct {
	ID        int `db:"id" json:"id"`
	UserID    int `db:"user_id" json:"user_id"`
	ProductID int `db:"product_id" json:"product_id"`
	Quantity  int `db:"quantity" json:"quantity"`
}

type DetailShoppingCarts []DetailShoppingCart
