package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	store, err := NewPostgresStorage()

	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(os.Getenv("SERVER_PORT"), store)
	server.Run()
	// fmt.Println("Hello World from MoneyGuardian")
}
