package handler

import (
	"encoding/json"
	"max-odyssey-backend/internal/service"
	"max-odyssey-backend/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ExamResultHandler struct {
	examResultService *service.ExamResultService
}

func NewExamResultHandler(examResultService *service.ExamResultService) *ExamResultHandler {
	return &ExamResultHandler{
		examResultService: examResultService,
	}
}

// GetUserExamResults handles GET /api/exam-results
func (h *ExamResultHandler) GetUserExamResults(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	results, err := h.examResultService.GetExamResultsByUserID(r.Context(), int32(userIDInt))
	if err != nil {
		http.Error(w, "Failed to fetch exam results", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, results)
}

// GetExamResultByID handles GET /api/exam-results/{id}
func (h *ExamResultHandler) GetExamResultByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid exam result ID", http.StatusBadRequest)
		return
	}

	result, err := h.examResultService.GetExamResultByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch exam result", http.StatusInternalServerError)
		return
	}

	// Check if the result belongs to the authenticated user
	userID := r.Context().Value("userID").(string)
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if int32(userIDInt) != result.UserID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}

// CreateExamResult handles POST /api/exam-results
func (h *ExamResultHandler) CreateExamResult(w http.ResponseWriter, r *http.Request) {
	var req service.CreateExamResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the user can only create exam results for themselves
	userID := r.Context().Value("userID").(string)
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Set the user ID from the authenticated user
	req.UserID = int64(userIDInt)

	result, err := h.examResultService.CreateExamResult(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create exam result", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, result)
}

// DeleteExamResult handles DELETE /api/exam-results/{id}
func (h *ExamResultHandler) DeleteExamResult(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid exam result ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.examResultService.DeleteExamResult(r.Context(), id, int32(userIDInt))
	if err != nil {
		http.Error(w, "Failed to delete exam result", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Exam result deleted successfully"})
}
