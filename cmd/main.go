package main

import (
	"log"
	"max-odyssey-backend/internal/database"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/handler"
	"max-odyssey-backend/internal/middleware"
	"net/http"

	"max-odyssey-backend/internal/service"

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
	r.Get("/questions", questionHandler.GetQuestions)
	//  curl -X GET http://localhost:8080/questions | jq '.' use and see what u get
	// docker exec -i max_odyssey-backend-postgres-1 psql -U satapp -d sat_tracker < sql/schema/004_add_choices_array.sql
	// running migrations
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
