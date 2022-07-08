package api

import (
	"companies/controllers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	server.Initialize(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	server.Run()
}
