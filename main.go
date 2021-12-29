package main

import (
	"log"

	"github.com/Imag/hyper-cache-api/app"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Could not find .env file")
	}

	server := app.New()
	server.Run()
}