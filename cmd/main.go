package main

import (
	"log"
	"max-odyssey-backend/internal/database"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/handler"
	"max-odyssey-backend/internal/middleware"
	"max-odyssey-backend/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	config := &database.Config{
		Host:     "localhost",
		Port:     5431,
		User:     "satapp",
		Password: "satapp123",
		DBName:   "sat_tracker",
	}

	dbConn, err := database.NewConnection(config)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)
	questionService := service.NewQuestionService(queries)
	questionHandler := handler.NewQuestionHandler(questionService)

	r := chi.NewRouter()
	r.Use(middleware.Cors)
	r.Get("/maxsat/practice", questionHandler.GetQuestions)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
