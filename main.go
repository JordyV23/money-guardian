package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Intenta cargar las variables de entorno desde un archivo .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	// Crea una nueva instancia de almacenamiento PostgreSQL.
	store, err := NewPostgresStorage()

	if err != nil {
		log.Fatal(err)
	}

	// Inicializa el almacenamiento.
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Crea una nueva instancia de servidor de API utilizando el puerto obtenido de las variables de entorno
	// y el almacenamiento creado anteriormente.
	server := NewAPIServer(os.Getenv("SERVER_PORT"), store)

	// Inicia el servidor de la API.
	server.Run()
	// fmt.Println("Hello World from MoneyGuardian")
}
