package handler

import (
	"encoding/json"
	"log"
	"max-odyssey-backend/internal/middleware"
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
	// Try to get the user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		log.Printf("No authenticated user, returning empty results")
		utils.RespondWithJSON(w, http.StatusOK, []service.ExamResultResponse{})
		return
	}

	log.Printf("Fetching exam results for user ID: %d", user.ID)
	log.Printf("Handler: User from context: ID=%d, Username=%s", user.ID, user.Username)

	results, err := h.examResultService.GetExamResultsByUserID(r.Context(), int32(user.ID))
	if err != nil {
		log.Printf("Error fetching exam results: %v", err)
		http.Error(w, "Failed to fetch exam results", http.StatusInternalServerError)
		return
	}

	// If results is nil, return an empty array instead
	if results == nil {
		log.Printf("No exam results found, returning empty array")
		utils.RespondWithJSON(w, http.StatusOK, []service.ExamResultResponse{})
		return
	}

	log.Printf("Found %d exam results for user ID %d", len(results), user.ID)
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
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if int32(user.ID) != result.UserID {
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
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Set the user ID from the authenticated user
	req.UserID = int64(user.ID)

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

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err = h.examResultService.DeleteExamResult(r.Context(), id, int32(user.ID))
	if err != nil {
		http.Error(w, "Failed to delete exam result", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Exam result deleted successfully"})
}
