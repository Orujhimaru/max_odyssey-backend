package handler

import (
	"encoding/json"
	"log"
	"max-odyssey-backend/internal/middleware"
	"max-odyssey-backend/internal/service"
	"max-odyssey-backend/utils"
	"net/http"
	"strconv"
)

type UserQuestionHandler struct {
	service *service.UserQuestionService
}

func NewUserQuestionHandler(service *service.UserQuestionService) *UserQuestionHandler {
	return &UserQuestionHandler{service: service}
}

func (h *UserQuestionHandler) GetBookmarkedQuestions(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Println("User not found in context")
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get sort direction from query parameters
	sortDir := r.URL.Query().Get("sort_dir")
	// log.Println(sortDir, "oruj")
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc" // Default to ascending if not specified or invalid
	}

	log.Printf("Getting bookmarks for user: %s (ID: %d) with sort_dir: %s", user.Username, user.ID, sortDir)

	// Get bookmarked questions
	questions, totalCount, err := h.service.GetBookmarkedQuestions(int64(user.ID), sortDir)
	if err != nil {
		log.Printf("Error getting bookmarks: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch bookmarked questions", err)
		return
	}

	log.Printf("Found %d bookmarked questions", len(questions))
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"questions":   questions,
		"total_count": totalCount,
		"sorting": map[string]string{
			"sort_dir": sortDir,
		},
	})
}

func (h *UserQuestionHandler) ToggleBookmark(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Println("User not found in context")
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	log.Printf("Toggle bookmark for user: %s (ID: %d)", user.Username, user.ID)

	// Log all query parameters
	log.Printf("Query parameters: %v", r.URL.Query())

	// Try to get question_id from both query and body
	questionIDStr := r.URL.Query().Get("question_id")
	log.Printf("question_id from query: %q", questionIDStr)

	// If not in query, try to get from body
	if questionIDStr == "" {
		var requestBody struct {
			QuestionID int64 `json:"question_id"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if requestBody.QuestionID == 0 {
			log.Println("No question_id found in body")
			utils.RespondWithError(w, http.StatusBadRequest, "Missing question_id", nil)
			return
		}

		log.Printf("question_id from body: %d", requestBody.QuestionID)
		questionID := requestBody.QuestionID

		// Toggle bookmark
		err = h.service.ToggleBookmark(int64(user.ID), questionID)
		if err != nil {
			log.Printf("Error toggling bookmark: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle bookmark", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	// Parse question ID from query
	questionID, err := strconv.ParseInt(questionIDStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing question_id: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid question ID", err)
		return
	}

	// Toggle bookmark
	err = h.service.ToggleBookmark(int64(user.ID), questionID)
	if err != nil {
		log.Printf("Error toggling bookmark: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle bookmark", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func (h *UserQuestionHandler) ToggleSolved(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Println("User not found in context")
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	log.Printf("Toggle solved for user: %s (ID: %d)", user.Username, user.ID)

	// Log all query parameters
	log.Printf("Query parameters: %v", r.URL.Query())

	// Try to get question_id from both query and body
	questionIDStr := r.URL.Query().Get("question_id")
	log.Printf("question_id from query: %q", questionIDStr)

	// If not in query, try to get from body
	if questionIDStr == "" {
		var requestBody struct {
			QuestionID int64 `json:"question_id"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if requestBody.QuestionID == 0 {
			log.Println("No question_id found in body")
			utils.RespondWithError(w, http.StatusBadRequest, "Missing question_id", nil)
			return
		}

		log.Printf("question_id from body: %d", requestBody.QuestionID)
		questionID := requestBody.QuestionID

		// Toggle solved status
		err = h.service.ToggleSolved(int64(user.ID), questionID)
		if err != nil {
			log.Printf("Error toggling solved status: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle solved status", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})
		return
	}

	// Parse question ID from query
	questionID, err := strconv.ParseInt(questionIDStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing question_id: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid question ID", err)
		return
	}

	// Toggle solved status
	err = h.service.ToggleSolved(int64(user.ID), questionID)
	if err != nil {
		log.Printf("Error toggling solved status: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle solved status", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
