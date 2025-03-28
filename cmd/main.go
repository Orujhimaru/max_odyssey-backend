package main

import (
	"encoding/json"
	"log"
	"max-odyssey-backend/internal/database"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/handler"
	"max-odyssey-backend/internal/middleware"
	"max-odyssey-backend/internal/service"
	"net/http"
	"time"

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

	// Initialize services
	questionService := service.NewQuestionService(queries)
	// userService := service.NewUserService(queries)  // Comment out or remove this line
	userQuestionService := service.NewUserQuestionService(queries)
	userSkillService := service.NewUserSkillService(queries)
	examResultService := service.NewExamResultService(queries)
	// JWT configuration
	jwtSecret := "your-secret-key" // In production, use environment variable
	jwtExpires := 24 * time.Hour
	authService := service.NewAuthService(queries, jwtSecret, jwtExpires)

	// Initialize handlers
	questionHandler := handler.NewQuestionHandler(questionService)
	userQuestionHandler := handler.NewUserQuestionHandler(userQuestionService)
	userSkillHandler := handler.NewUserSkillHandler(userSkillService)
	examResultHandler := handler.NewExamResultHandler(examResultService)
	// Create auth middleware
	authMiddleware := middleware.Auth(authService)

	r := chi.NewRouter()
	r.Use(middleware.Cors)

	// Public routes - order matters!
	r.Get("/questions/filtered", questionHandler.GetFilteredQuestions)
	r.Get("/questions/{id}", questionHandler.GetQuestion)
	r.Get("/questions", questionHandler.GetQuestions)
	r.Get("/exams", examResultHandler.GetUserExamResults)
	// r.Get("/test-direct-query", questionHandler.TestDirectQuery)
	//  curl -X GET http://localhost:8080/questions | jq '.' use and see what u get
	// docker exec -i max_odyssey-backend-postgres-1 psql -U satapp -d sat_tracker < sql/schema/004_add_choices_array.sql
	// running migrations
	// Auth routes (login, register, etc.)
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Login request received")

		// Parse the request body
		var loginReq struct {
			Username string `json:"username"`
		}

		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			log.Printf("Error parsing login request: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("Login attempt for user: %s", loginReq.Username)

		// Get token for the user
		token, err := authService.Login(loginReq.Username)
		if err != nil {
			log.Printf("Login error: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("Login successful for user: %s", loginReq.Username)

		// Return the token
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token":"` + token + `"}`))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		// User question routes
		r.Get("/bookmarks", userQuestionHandler.GetBookmarkedQuestions)
		r.Post("/bookmark", userQuestionHandler.ToggleBookmark)
		r.Post("/solved", userQuestionHandler.ToggleSolved)

		// User skill routes
		r.Get("/skills", userSkillHandler.GetUserSkills)
		r.Post("/skills", userSkillHandler.CreateOrUpdateUserSkill)
		r.Delete("/skills", userSkillHandler.DeleteUserSkill)

		// Add more protected routes here
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
