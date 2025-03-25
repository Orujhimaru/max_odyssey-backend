package handler

import (
	"database/sql"
	"log"
	"max-odyssey-backend/internal/models"
	"max-odyssey-backend/internal/service"
	"max-odyssey-backend/utils"
	"net/http"
	"strconv"
	"strings"

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

// GetFilteredQuestions handles requests for filtered questions
func (h *QuestionHandler) GetFilteredQuestions(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("Filtered questions request received: %s", r.URL.String())

	// Parse query parameters
	subjectStr := r.URL.Query().Get("subject")
	difficultyStr := r.URL.Query().Get("difficulty")
	topicStr := r.URL.Query().Get("topic")
	subtopicStr := r.URL.Query().Get("subtopic")
	log.Printf("Extracted subtopic: %q", subtopicStr)
	sortDirStr := r.URL.Query().Get("sort_dir")
	pageSizeStr := r.URL.Query().Get("page_size")
	pageNumberStr := r.URL.Query().Get("page")

	// Log the parsed parameters
	log.Printf("Request parameters - Subject: %s, Difficulty: %s, Topic: %s, Subtopic: %s, SortDir: %s, PageSize: %s, PageNumber: %s",
		subjectStr, difficultyStr, topicStr, subtopicStr, sortDirStr, pageSizeStr, pageNumberStr)

	// Create filters struct
	filters := service.QuestionFilters{
		SortDir: sortDirStr,
	}

	// Parse subject (required)
	subject, err := strconv.Atoi(subjectStr)
	if err != nil {
		log.Printf("Error parsing subject ID: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid subject ID", err)
		return
	}
	filters.SubjectID = subject

	// Parse difficulty
	if difficultyStr != "" {
		difficulty, err := strconv.Atoi(difficultyStr)
		if err == nil && difficulty >= 0 && difficulty <= 2 {
			filters.DifficultyLevel = difficulty
			log.Printf("Set difficulty level to: %d", filters.DifficultyLevel)
		} else {
			log.Printf("Invalid difficulty level: %s, ignoring", difficultyStr)
			filters.DifficultyLevel = -1 // Special value for "no filter"
		}
	} else {
		filters.DifficultyLevel = -1 // Special value for "no filter"
	}

	// Parse topic
	if topicStr != "" {
		filters.Topic = topicStr
	}

	// Parse subtopic
	if subtopicStr != "" {
		// Remove commas from subtopic
		trimmedSubtopic := strings.ReplaceAll(subtopicStr, ",", "")
		log.Printf("Original subtopic: %q, Trimmed subtopic: %q", subtopicStr, trimmedSubtopic)

		filters.Subtopic = trimmedSubtopic
	}

	// Parse page size
	if pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSize > 0 && pageSize <= 100 {
			filters.PageSize = pageSize
		} else {
			log.Printf("Invalid page size: %s, using default", pageSizeStr)
			filters.PageSize = 10 // Default
		}
	} else {
		filters.PageSize = 10 // Default
	}

	// Parse page number
	if pageNumberStr != "" {
		pageNumber, err := strconv.Atoi(pageNumberStr)
		if err == nil && pageNumber > 0 {
			filters.PageNumber = pageNumber
		} else {
			log.Printf("Invalid page number: %s, using default", pageNumberStr)
			filters.PageNumber = 1 // Default
		}
	} else {
		filters.PageNumber = 1 // Default
	}

	// Log the processed filters
	log.Printf("Processed filters - Subject: %d, Difficulty: %v, Topic: %v, Subtopic: %v, SortDir: %s, PageSize: %d, PageNumber: %d",
		filters.SubjectID, filters.DifficultyLevel, filters.Topic, filters.Subtopic, filters.SortDir, filters.PageSize, filters.PageNumber)

	// Get filtered questions
	questions, totalCount, err := h.service.GetFilteredQuestions(filters)
	if err != nil {
		log.Printf("Error fetching filtered questions: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch questions", err)
		return
	}

	// Log the query results
	log.Printf("Query results - Total count: %d, Questions returned: %d", totalCount, len(questions))

	// Calculate total pages
	totalPages := (totalCount + filters.PageSize - 1) / filters.PageSize

	// Prepare response
	response := struct {
		Questions  []models.Question `json:"questions"`
		Pagination struct {
			CurrentPage int `json:"current_page"`
			PageSize    int `json:"page_size"`
			TotalItems  int `json:"total_items"`
			TotalPages  int `json:"total_pages"`
		} `json:"pagination"`
		Filters struct {
			Subject    int    `json:"subject"`
			Difficulty int    `json:"difficulty,omitempty"`
			Topic      string `json:"topic,omitempty"`
			Subtopic   string `json:"subtopic,omitempty"`
		} `json:"filters"`
		Sorting struct {
			SortDir string `json:"sort_dir"`
		} `json:"sorting"`
	}{
		Questions: questions,
		Pagination: struct {
			CurrentPage int `json:"current_page"`
			PageSize    int `json:"page_size"`
			TotalItems  int `json:"total_items"`
			TotalPages  int `json:"total_pages"`
		}{
			CurrentPage: filters.PageNumber,
			PageSize:    filters.PageSize,
			TotalItems:  totalCount,
			TotalPages:  totalPages,
		},
		Filters: struct {
			Subject    int    `json:"subject"`
			Difficulty int    `json:"difficulty,omitempty"`
			Topic      string `json:"topic,omitempty"`
			Subtopic   string `json:"subtopic,omitempty"`
		}{
			Subject:    filters.SubjectID,
			Difficulty: filters.DifficultyLevel,
			Topic:      filters.Topic,
			Subtopic:   filters.Subtopic,
		},
		Sorting: struct {
			SortDir string `json:"sort_dir"`
		}{
			SortDir: filters.SortDir,
		},
	}

	// Log a sample of the response (first question if available)
	if len(questions) > 0 {
		firstQuestion := questions[0]
		log.Printf("Sample question - ID: %d, Text: %.50s..., Topic: %s, Difficulty: %d, SolveRate: %d",
			firstQuestion.ID, firstQuestion.QuestionText, firstQuestion.Topic, firstQuestion.DifficultyLevel, firstQuestion.SolveRate)
	}

	// Log pagination details
	log.Printf("Pagination - CurrentPage: %d, PageSize: %d, TotalItems: %d, TotalPages: %d",
		filters.PageNumber, filters.PageSize, totalCount, totalPages)

	utils.RespondWithJSON(w, http.StatusOK, response)
	log.Printf("Response sent successfully")
}
