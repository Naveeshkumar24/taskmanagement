package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/naveeshkumar24/pkg/database"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load env: %v", err)
	}
	conn := NewConnection()
	defer conn.DB.Close()
	server := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: registerTaskRouter(conn.DB),
	}
	query := database.NewQuery(conn.DB)
	err := query.CreateTaskTables()
	if err != nil {
		log.Fatal("Unable to create database %v", err)
	}

	log.Printf("server is running at port %s", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start the server: %v", err)
	}
}
