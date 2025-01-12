package models

import "time"

type User struct {
	Id         int        `db:"id" json:"id"`
	Username   string     `db:"username" json:"username"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"password"`
	Created_at *time.Time `db:"created_at" json:"-"`
	Updated_at *time.Time `db:"updated_at" json:"-"`
}

type Users []User
