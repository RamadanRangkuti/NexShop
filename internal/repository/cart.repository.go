package repository

import (
	"database/sql"

	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type ShoppingCartRepositoryInterface interface {
	FindCartItem(userID, productID int) (*models.ShoppingCart, error)
	AddCartItem(userID, productID, quantity int) error
	UpdateCartItem(userID, productID, quantity int) error
	FindCartByUserid(id int) (*models.DetailShoppingCarts, error)
	ClearCartByUserID(userID int) error
}

type ShoppingCartRepository struct {
	*sqlx.DB
}

func NewShoppingCartRepository(db *sqlx.DB) *ShoppingCartRepository {
	return &ShoppingCartRepository{db}
}

func (r *ShoppingCartRepository) FindCartItem(userID, productID int) (*models.ShoppingCart, error) {
	query := `SELECT id, user_id, product_id, quantity FROM shopping_cart 
	          WHERE user_id = $1 AND product_id = $2`
	data := models.ShoppingCart{}
	err := r.Get(&data, query, userID, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *ShoppingCartRepository) AddCartItem(userID, productID, quantity int) error {
	query := `INSERT INTO shopping_cart (user_id, product_id, quantity) 
	          VALUES ($1, $2, $3)`
	_, err := r.Exec(query, userID, productID, quantity)
	return err
}

func (r *ShoppingCartRepository) UpdateCartItem(userID, productID, quantity int) error {
	query := `UPDATE shopping_cart 
	          SET quantity = quantity + $1 
	          WHERE user_id = $2 AND product_id = $3`
	_, err := r.Exec(query, quantity, userID, productID)
	return err
}

func (r *ShoppingCartRepository) FindCartByUserid(id int) (*models.DetailShoppingCarts, error) {
	query := `select sc.id, u.username, p.name, sc.quantity from users u join shopping_cart sc 
on u.id = sc.user_id join products p on p.id = sc.product_id where user_id=$1`
	data := models.DetailShoppingCarts{}

	err := r.Select(&data, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *ShoppingCartRepository) ClearCartByUserID(userID int) error {
	query := `DELETE FROM shopping_cart WHERE user_id = $1`
	_, err := r.Exec(query, userID)
	return err
}
