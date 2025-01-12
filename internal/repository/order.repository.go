package repository

import (
	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrderRepositoryInterface interface {
	CreateOrder(userID int, total float64) (*models.Order, error)
	AddOrderItem(orderID, productID, quantity int, price float64) error
}

type OrderRepository struct {
	*sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) CreateOrder(userID int, total float64) (*models.Order, error) {
	query := `INSERT INTO orders (customer_id, amount, order_date) 
	          VALUES ($1, $2,CURRENT_TIMESTAMP) RETURNING id, customer_id, amount, order_date`
	order := &models.Order{}
	err := r.Get(order, query, userID, total)
	return order, err
}

func (r *OrderRepository) AddOrderItem(orderID, productID, quantity int, price float64) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price) 
	          VALUES ($1, $2, $3, $4)`
	_, err := r.Exec(query, orderID, productID, quantity, price)
	return err
}
