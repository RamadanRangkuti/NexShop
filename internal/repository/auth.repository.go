package repository

import (
	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepositoryInterface interface {
	RegisterUser(body *models.Auth) (*models.Auth, error)
}

type AuthRepository struct {
	*sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) RegisterUser(body *models.Auth) (*models.Auth, error) {
	query := `INSERT INTO users (
		email,
		username,
		password
	  ) VALUES (
		:email,
		:username,
		:password
	  ) RETURNING id, email, password`

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
