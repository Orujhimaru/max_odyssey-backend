package handler

import (
	"max-odyssey-backend/internal/models"
	"max-odyssey-backend/internal/service"
	"max-odyssey-backend/utils"
	"net/http"
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
