package main

import (
	"math/rand"
	"time"
)

// CreateAccountRequest es una estructura utilizada para representar una solicitud de creación de cuenta en formato JSON.
type CreateAccountRequest struct {
	FirstName string `json:"firstName"` // El campo FirstName almacena el primer nombre del titular de la cuenta.
	LastName  string `json:"lastName"`  // El campo LastName almacena el apellido del titular de la cuenta.
}

// Account es una estructura que representa una cuenta.
type Account struct {
	ID        int       `json:"id"`            // El campo ID almacena el identificador único de la cuenta.
	Firstname string    `json:"firstName"`     // El campo Firstname almacena el nombre del titular de la cuenta.
	Lastname  string    `json:"lastName"`      // El campo Lastname almacena el apellido del titular de la cuenta.
	Number    int64     `json:"accountNumber"` // El campo Number almacena el número de cuenta, que puede ser un número largo.
	Balance   int64     `json:"balance"`       // El campo Balance almacena el saldo de la cuenta.
	CreatedAt time.Time `json:"createdAt"`     // El campo CreatedAt almacena fecha de creacion de la cuenta.
}

// NewAccount crea una nueva instancia de Account y la inicializa con los datos proporcionados.
// Parámetros:
// - firstname: El nombre del titular de la cuenta.
// - lastname: El apellido del titular de la cuenta.
// Devuelve:
// - Un puntero a la instancia de Account creada con valores aleatorios para ID y Number.
func NewAccount(firstname, lastname string) *Account {
	// Crea una nueva instancia de Account utilizando un literal compuesto.
	// Asigna valores aleatorios para ID y Number.
	return &Account{
		Firstname: firstname,                 // Asigna el nombre del titular de la cuenta.
		Lastname:  lastname,                  // Asigna el apellido del titular de la cuenta.
		Number:    int64(rand.Intn(1000000)), // Asigna un número de cuenta aleatorio en el rango [0, 999999].
		CreatedAt: time.Now().UTC(),          // Asigna la fecha y hora actual como fecha de creación de la cuenta.
	}
}
