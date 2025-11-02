package main

import (
	"log"

	app "github.com/IlyaChgn/voblako/internal/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Println(".env file not found, using OS environment")
	}

	srv := new(app.Server)
	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server ", err)
	}
}
