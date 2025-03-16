package service

import (
	"context"
	"log"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/models"
)

type QuestionService struct {
	db *db.Queries
}

func NewQuestionService(db *db.Queries) *QuestionService {
	return &QuestionService{
		db: db,
	}
}

func (s *QuestionService) GetQuestions() ([]models.Question, error) {
	ctx := context.Background()
	dbQuestions, err := s.db.GetQuestions(ctx)
	if err != nil {
		return nil, err
	}

	// log.Printf("Got %d questions from database", len(dbQuestions))
	// if len(dbQuestions) > 0 {
	// 	log.Printf("First question choices: %v", dbQuestions[0].Choices)
	// }

	questions := make([]models.Question, len(dbQuestions))
	for i, q := range dbQuestions {
		questions[i] = models.Question{
			ID:                 int(q.ID),
			SubjectID:          int(q.SubjectID.Int32),
			QuestionText:       q.QuestionText,
			Topic:              q.Topic.String,
			Subtopic:           q.Subtopic.String,
			CorrectAnswerIndex: int(q.CorrectAnswerIndex.Int32),
			DifficultyLevel:    int(q.DifficultyLevel.Int32),
			Explanation:        q.Explanation.String,
			CreatedAt:          q.CreatedAt.Time,
			SolveRate:          int(q.SolveRate.Int32),
			Choices:            q.Choices, // Use the actual choices from DB, not an empty array
		}
	}

	return questions, nil
}

// func (s *QuestionService) CreateQuestion(question db.CreateQuestionParams) error {
// 	return s.db.CreateQuestion(context.Background(), question)
// }

// GetQuestionByID gets a single question by ID
func (s *QuestionService) GetQuestionByID(id int64) (*models.Question, error) {
	q, err := s.db.GetQuestion(context.Background(), int32(id))
	if err != nil {
		return nil, err
	}

	return &models.Question{
		ID:                 int(q.ID),
		SubjectID:          int(q.SubjectID.Int32),
		QuestionText:       q.QuestionText,
		Topic:              q.Topic.String,
		Subtopic:           q.Subtopic.String,
		CorrectAnswerIndex: int(q.CorrectAnswerIndex.Int32),
		DifficultyLevel:    int(q.DifficultyLevel.Int32),
		Explanation:        q.Explanation.String,
		CreatedAt:          q.CreatedAt.Time,
		SolveRate:          int(q.SolveRate.Int32),
		Choices:            q.Choices, // Just use the choices directly
	}, nil
}

// QuestionFilters represents all possible filters for questions
type QuestionFilters struct {
	SubjectID       int
	DifficultyLevel int
	Topic           string
	Subtopic        string
	SortDir         string
	PageSize        int
	PageNumber      int
}

// GetFilteredQuestions gets questions with filtering, sorting and pagination
func (s *QuestionService) GetFilteredQuestions(filters QuestionFilters) ([]models.Question, int, error) {
	ctx := context.Background()

	// Set defaults
	if filters.PageNumber <= 0 {
		filters.PageNumber = 1
	}

	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}

	if filters.SortDir == "" {
		filters.SortDir = "asc"
	}

	// Calculate offset
	offset := (filters.PageNumber - 1) * filters.PageSize

	// Prepare parameters with special values for "no filter"
	difficultyParam := int32(filters.DifficultyLevel)
	if filters.DifficultyLevel < 0 || filters.DifficultyLevel > 2 {
		difficultyParam = -1 // Special value for "no filter"
	}

	// Log the parameters
	log.Printf("Executing GetFilteredQuestions with filters: %+v", filters)
	log.Printf("SQL Parameters - SubjectID: %d, DifficultyLevel: %d, Topic: %s, Subtopic: %s, SortDir: %s, PageSize: %d, PageOffset: %d",
		filters.SubjectID,
		difficultyParam,
		filters.Topic,
		filters.Subtopic,
		filters.SortDir,
		filters.PageSize,
		offset)

	// Get questions with total count
	dbQuestions, err := s.db.GetFilteredQuestions(ctx, db.GetFilteredQuestionsParams{
		SubjectID: filters.SubjectID,
		Column2:   difficultyParam,
		Column3:   filters.Topic,
		Column4:   filters.Subtopic,
		Column5:   filters.SortDir,
		Limit:     int32(filters.PageSize),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	// Extract total count from first row
	var totalCount int
	if len(dbQuestions) > 0 {
		totalCount = int(dbQuestions[0].TotalCount)
	}

	// Convert to models
	questions := make([]models.Question, len(dbQuestions))
	for i, q := range dbQuestions {
		questions[i] = models.Question{
			ID:                 int(q.ID),
			SubjectID:          int(q.SubjectID.Int32),
			QuestionText:       q.QuestionText,
			Topic:              q.Topic.String,
			Subtopic:           q.Subtopic.String,
			CorrectAnswerIndex: int(q.CorrectAnswerIndex.Int32),
			DifficultyLevel:    int(q.DifficultyLevel.Int32),
			Explanation:        q.Explanation.String,
			SolveRate:          int(q.SolveRate.Int32),
			Choices:            q.Choices,
		}
	}

	return questions, totalCount, nil
}
