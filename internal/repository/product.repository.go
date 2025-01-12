package repository

import (
	"database/sql"
	"fmt"

	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepositoryInterface interface {
	FindAllProduct(limit, offset int, search, sort, order string) (*models.Products, error)
	FindProductById(id int) (*models.Product, error)
	InsertProduct(body *models.Product) (*models.Product, error)
	EditProduct(id int, body *models.Product) (*models.Product, error)
	RemoveProduct(id int) error
	CountProduct(search string) (int, error)
	ReduceStock(productID, quantity int) error
}

type ProductRepository struct {
	*sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) FindAllProduct(limit, offset int, search, sort, order string) (*models.Products, error) {
	if sort != "price" {
		sort = "id"
	}

	if order != "desc" {
		order = "asc"
	}

	query := `SELECT id, name, description, price, stock 
	          FROM products
	          WHERE LOWER(name) LIKE $1
	          ORDER BY ` + sort + ` ` + order + `
	          LIMIT $2 OFFSET $3`

	data := models.Products{}
	err := r.Select(&data, query, "%"+search+"%", limit, offset)
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

func (r *ProductRepository) CountProduct(search string) (int, error) {
	query := `SELECT COUNT(*) FROM products WHERE LOWER(name) LIKE $1`
	var count int
	err := r.Get(&count, query, "%"+search+"%")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductRepository) ReduceStock(productID, quantity int) error {
	query := `UPDATE products 
	          SET stock = stock - $1 
	          WHERE id = $2 AND stock >= $1`
	result, err := r.Exec(query, quantity, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("insufficient stock")
	}

	return nil
}
