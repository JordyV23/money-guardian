package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=" + os.Getenv("DB_USER") + " dbname=" + os.Getenv("DB_NAME") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) CreateAccount(*Account) error {
	return nil
}
func (s *PostgresStorage) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStorage) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStorage) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
