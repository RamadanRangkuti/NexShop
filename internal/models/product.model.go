package models

import "time"

type Product struct {
	Id          int        `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	Price       float64    `db:"price" json:"price"`
	Stock       int        `db:"stock" json:"stock"`
	Created_at  *time.Time `db:"created_at" json:"-"`
	Updated_at  *time.Time `db:"updated_at" json:"-"`
}

type Products []Product
