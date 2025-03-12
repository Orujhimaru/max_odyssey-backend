package main

import (
	"log"
	"max-odyssey-backend/internal/database"
)

func main() {
	config := &database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "satapp",
		Password: "satapp123",
		DBName:   "sat_tracker",
	}

	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")
}
