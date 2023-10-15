package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Interfaz que define un conjunto de métodos para realizar operaciones en un sistema de almacenamiento.
type Storage interface {
	// CreateAccount crea una nueva cuenta y almacena los datos en el sistema de almacenamiento.
	CreateAccount(*Account) error

	// DeleteAccount elimina una cuenta del sistema de almacenamiento basándose en su identificador único.
	DeleteAccount(int) error

	// UpdateAccount actualiza los datos de una cuenta existente en el sistema de almacenamiento.
	UpdateAccount(*Account) error

	// GetAccounts recupera todas las cuentas almacenadas en el sistema de almacenamiento.
	GetAccounts() ([]*Account, error)

	// GetAccountById recupera una cuenta basada en su identificador único desde el sistema de almacenamiento.
	GetAccountById(int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage crea una nueva instancia de PostgresStorage para interactuar con una base de datos PostgreSQL.
// Parámetros:
// - dbname: El nombre de la base de datos PostgreSQL.
// - user: El nombre de usuario de la base de datos PostgreSQL.
// - password: La contraseña para acceder a la base de datos PostgreSQL.
// Devuelve:
// - Un puntero a la instancia de PostgresStorage creada y una posible instancia de error en caso de problemas durante la creación.
func NewPostgresStorage() (*PostgresStorage, error) {
	// Construye la cadena de conexión a la base de datos PostgreSQL utilizando valores de las variables de entorno.
	connStr := "user=" + os.Getenv("DB_USER") + " dbname=" + os.Getenv("DB_NAME") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"

	// Abre una conexión a la base de datos PostgreSQL.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verifica si la conexión a la base de datos es exitosa.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Crea una instancia de PostgresStorage y la inicializa con la conexión a la base de datos.
	return &PostgresStorage{
		db: db,
	}, nil
}

// Init inicializa el sistema de almacenamiento PostgreSQL, creando la tabla "account" si aún no existe.
// Devuelve:
// - Un error en caso de que ocurra algún problema durante la inicialización, de lo contrario, devuelve nil.
func (s *PostgresStorage) Init() error {
	// Llama al método createAccountTable para crear la tabla "account" en la base de datos si aún no existe.
	return s.createAccountTable()
}

// createAccountTable crea una tabla "account" en la base de datos PostgreSQL si aún no existe.
// Devuelve:
// - Un error en caso de que ocurra algún problema durante la creación de la tabla, de lo contrario, devuelve nil.
func (s *PostgresStorage) createAccountTable() error {
	// Definición de la consulta SQL para crear la tabla "account".
	query := `CREATE TABLE IF NOT EXISTS account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		create_at timestamp default current_timestamp
	)`

	// Ejecuta la consulta en la base de datos PostgreSQL.
	_, err := s.db.Exec(query)

	// Devuelve un posible error si ocurrió algún problema durante la creación de la tabla.
	return err
}

func (s *PostgresStorage) CreateAccount(account *Account) error {
	query := `
	INSERT INTO account 
	(
		first_name,
		last_name,
		number,
		balance,
		create_at
	) VALUES  (
		$1,
		$2,
		$3,
		$4,
		$5
	)`

	resp, err := s.db.Exec(
		query,
		account.Firstname,
		account.Lastname,
		account.Number,
		account.Balance,
		account.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}
func (s *PostgresStorage) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}
func (s *PostgresStorage) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT id,first_name, last_name, number, balance, create_at FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("no se encontro la cuenta con el id %d", id)
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT id,first_name, last_name, number, balance, create_at FROM account")

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.Firstname,
		&account.Lastname,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	return account, err

}
