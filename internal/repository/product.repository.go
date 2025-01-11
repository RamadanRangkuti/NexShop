package repository

import (
	"database/sql"

	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepositoryInterface interface {
	FindAllProduct() (*models.Products, error)
	FindProductById(id int) (*models.Product, error)
	InsertProduct(body *models.Product) (*models.Product, error)
	EditProduct(id int, body *models.Product) (*models.Product, error)
	RemoveProduct(id int) error
}

type ProductRepository struct {
	*sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) FindAllProduct() (*models.Products, error) {
	query := `SELECT id, name, description, price, stock  FROM products ORDER BY id asc`
	data := models.Products{}

	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *ProductRepository) FindProductById(id int) (*models.Product, error) {
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id=$1`
	data := models.Product{}

	err := r.Get(&data, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *ProductRepository) InsertProduct(body *models.Product) (*models.Product, error) {
	query := `INSERT INTO products (
                name,
                description,
                price,
                stock
              ) VALUES (
                :name,
                :description,
                :price,
                :stock
              ) RETURNING id, name, description, price, stock, created_at`

	rows, err := r.NamedQuery(query, body)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(body)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func (r *ProductRepository) EditProduct(id int, body *models.Product) (*models.Product, error) {
	query := `UPDATE products 
              SET name = :name, 
                  description = :description, 
                  price = :price, 
                  stock = :stock,
				  updated_at = :updated_at 
              WHERE id = :id`

	params := map[string]interface{}{
		"id":          id,
		"name":        body.Name,
		"description": body.Description,
		"price":       body.Price,
		"stock":       body.Stock,
		"updated_at":  body.Updated_at,
	}

	_, err := r.NamedExec(query, params)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *ProductRepository) RemoveProduct(id int) error {
	query := `DELETE FROM products WHERE id=$1`
	_, err := r.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
