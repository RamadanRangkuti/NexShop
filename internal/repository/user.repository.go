package repository

import (
	"database/sql"

	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryInterface interface {
	FindAllUser(limit, offset int, search, sort, order string) (*models.Users, error)
	FindUserById(id int) (*models.User, error)
	InsertUser(body *models.User) (*models.User, error)
	EditUser(id int, body *models.User) (*models.User, error)
	RemoveUser(id int) error
	CountUser(search string) (int, error)

	FindUserByUsername(username string) (*models.User, error)
	FindUsersBySignupDate(date string) (*models.Users, error)
	FindUserByEmail(email string) (*models.User, error)
}

type UserRepository struct {
	*sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) FindAllUser(limit, offset int, search, sort, order string) (*models.Users, error) {
	if sort != "price" {
		sort = "id"
	}

	if order != "desc" {
		order = "asc"
	}

	query := `SELECT id, username, email, password 
	          FROM users
	          WHERE LOWER(username) LIKE $1
	          ORDER BY ` + sort + ` ` + order + `
	          LIMIT $2 OFFSET $3`

	data := models.Users{}
	err := r.Select(&data, query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserRepository) FindUserById(id int) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE id=$1`
	data := models.User{}

	err := r.Get(&data, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *UserRepository) InsertUser(body *models.User) (*models.User, error) {
	query := `INSERT INTO users (
                username,
                email,
                password
              ) VALUES (
                :username,
                :email,
                :password
              ) RETURNING id, username, email, password`

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

func (r *UserRepository) EditUser(id int, body *models.User) (*models.User, error) {
	query := `UPDATE users 
              SET username = :username, 
                  email = :email, 
                  password = :password,
				  updated_at = :updated_at 
              WHERE id = :id`

	params := map[string]interface{}{
		"id":         id,
		"username":   body.Username,
		"email":      body.Email,
		"password":   body.Password,
		"updated_at": body.Updated_at,
	}

	_, err := r.NamedExec(query, params)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *UserRepository) RemoveUser(id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	data := models.User{}
	err := r.Get(&data, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *UserRepository) FindUsersBySignupDate(date string) (*models.Users, error) {
	query := `SELECT id, username, email, password FROM users WHERE created_at > $1`
	data := models.Users{}
	err := r.Select(&data, query, date)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, COALESCE(username, '') AS username, email, password FROM users WHERE email = $1`
	data := models.User{}
	err := r.Get(&data, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *UserRepository) CountUser(search string) (int, error) {
	query := `SELECT COUNT(id) FROM users WHERE LOWER(username) LIKE $1`
	var count int
	err := r.Get(&count, query, "%"+search+"%")
	if err != nil {
		return 0, err
	}
	return count, nil
}
