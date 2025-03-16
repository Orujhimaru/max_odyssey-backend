package handler

import (
	"database/sql"
	"max-odyssey-backend/internal/models"
	"max-odyssey-backend/internal/service"
	"max-odyssey-backend/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type QuestionHandler struct {
	service *service.QuestionService
}

func NewQuestionHandler(service *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.service.GetQuestions()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch questions", err)
		return
	}

	type response struct {
		Questions []models.Question `json:"questions"`
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		Questions: questions,
	})
}

// GetQuestion gets a single question by ID
func (h *QuestionHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	// Get question ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid question ID", err)
		return
	}

	// Get the question
	question, err := h.service.GetQuestionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Question not found", nil)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch question", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, question)
}
