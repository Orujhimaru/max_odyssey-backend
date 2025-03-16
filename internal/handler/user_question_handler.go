package handler

import (
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

	log.Printf("Getting bookmarks for user: %s (ID: %d)", user.Username, user.ID)

	// Get bookmarked questions
	questions, err := h.service.GetBookmarkedQuestions(int64(user.ID))
	if err != nil {
		log.Printf("Error getting bookmarks: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch bookmarked questions", err)
		return
	}

	log.Printf("Found %d bookmarked questions", len(questions))
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"questions": questions,
	})
}

func (h *UserQuestionHandler) ToggleBookmark(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get question ID from request
	questionID, err := strconv.ParseInt(r.URL.Query().Get("question_id"), 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid question ID", err)
		return
	}

	// Toggle bookmark
	err = h.service.ToggleBookmark(int64(user.ID), questionID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle bookmark", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
