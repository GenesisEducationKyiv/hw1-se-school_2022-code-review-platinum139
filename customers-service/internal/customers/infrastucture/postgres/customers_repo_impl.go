package postgres

import (
	"context"
	"customers-service/internal/customers/domain"
	"database/sql"
	"github.com/lib/pq"
)

type CustomersRepoImpl struct {
	db *sql.DB
}

func (c *CustomersRepoImpl) CreateCustomer(customer domain.Customer) (*domain.Customer, error) {
	query := "INSERT INTO customers (email) VALUES ($1) RETURNING *;"

	row := c.db.QueryRow(query, customer.Email)
	if row.Err() != nil {
		if pgErr, ok := row.Err().(*pq.Error); ok {
			if pgErr.Code.Name() == "unique_violation" {
				return nil, domain.EmailAlreadyExistsError{}
			}
		}
		return nil, row.Err()
	}

	var createdCustomer domain.Customer
	if err := row.Scan(
		&createdCustomer.ID,
		&createdCustomer.Email,
	); err != nil {
		return nil, err
	}

	return &createdCustomer, nil
}

func (c *CustomersRepoImpl) DeleteCustomer(customerID int64) error {
	query := "DELETE FROM customers WHERE id = $1;"

	_, err := c.db.Exec(query, customerID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomersRepoImpl) CreateCustomerWithTransaction(
	customer domain.Customer, transaction domain.ProcessedTransaction) error {

	ctx := context.Background()
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := "INSERT INTO customers (email) VALUES ($1);"

	_, err = tx.Exec(query, customer.Email)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code.Name() == "unique_violation" {
				return domain.EmailAlreadyExistsError{}
			}
		}
		return err
	}

	query = "INSERT INTO transactions (transaction_id) VALUES ($1);"

	_, err = tx.Exec(query, transaction.TransactionID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *CustomersRepoImpl) DeleteCustomerWithTransaction(
	customer domain.Customer, transaction domain.ProcessedTransaction) error {

	ctx := context.Background()
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := "SELECT * FROM transactions WHERE transaction_id = $1;"

	row := tx.QueryRow(query, transaction.TransactionID)
	if row.Err() != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return row.Err()
	}
	var processedTransaction domain.ProcessedTransaction
	err = row.Scan(&processedTransaction.ID, &processedTransaction.TransactionID)
	if err == sql.ErrNoRows {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return nil
	}

	query = "DELETE FROM customers WHERE email = $1;"

	_, err = c.db.Exec(query, customer.Email)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewCustomersRepo(db *sql.DB) *CustomersRepoImpl {
	return &CustomersRepoImpl{
		db: db,
	}
}
