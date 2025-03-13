package handler

import (
	"encoding/json"
	"max-odyssey-backend/internal/service"
	"net/http"
)

type QuestionHandler struct {
	service *service.QuestionService
}

func NewQuestionHandler(service *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	questions, err := h.service.GetQuestions(ctx)
	if err != nil {
		http.Error(w, "Failed to get questions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}
