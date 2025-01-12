/*
Dengan transaksi database, kita memastikan thread-safety dan menghindari race conditions selama operasi deposit dan withdraw
*/
package repository

import (
	"database/sql"
	"errors"

	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/jmoiron/sqlx"
)

type AccountRepositoryInterface interface {
	Deposit(userID int, amount float64) error
	Withdraw(userID int, amount float64) error
	GetBalance(userID int) (float64, error)
	CreateAccount(userID int) error
}

type AccountRepository struct {
	*sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) Deposit(userID int, amount float64) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Update balance akun
	query := `UPDATE accounts SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2`
	result, err := tx.Exec(query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}

func (r *AccountRepository) Withdraw(userID int, amount float64) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Check apakah balance memadai
	var balance float64
	query := `SELECT balance FROM accounts WHERE user_id = $1`
	err = tx.Get(&balance, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("account not found")
		}
		return err
	}

	if balance < amount {
		return errors.New("insufficient balance")
	}

	// Update akun balance
	updateQuery := `UPDATE accounts SET balance = balance - $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2`
	result, err := tx.Exec(updateQuery, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}

func (r *AccountRepository) GetBalance(userID int) (float64, error) {
	var balance float64
	query := `SELECT balance FROM accounts WHERE user_id = $1`
	err := r.Get(&balance, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("account not found")
		}
		return 0, err
	}
	return balance, nil
}

func (r *AccountRepository) CreateAccount(userID int) error {
	query := `INSERT INTO accounts (
		user_id,
		balance
	  ) VALUES (
		:user_id,
		:balance
	  )`

	account := models.CreateAccount{
		UserId:  userID,
		Balance: 0,
	}

	_, err := r.NamedExec(query, account)
	return err
}
