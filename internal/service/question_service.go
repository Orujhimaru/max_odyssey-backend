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
			Choices:            q.Choices,
			Passage:            q.Passage.String,
			Bluebook:           q.Bluebook,
			HTMLTable:          q.HtmlTable.String,
			SVGImage:           q.SvgImage.String,
			IsMultipleChoice:   q.IsMultipleChoice.Bool,
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
		Choices:            q.Choices,
		Passage:            q.Passage.String,
		Bluebook:           q.Bluebook,
		HTMLTable:          q.HtmlTable.String,
		SVGImage:           q.SvgImage.String,
		IsMultipleChoice:   q.IsMultipleChoice.Bool,
	}, nil
}

// QuestionFilters represents all possible filters for questions
type QuestionFilters struct {
	SubjectID          int
	DifficultyLevel    int
	Topic              string
	Subtopic           string
	HasSpecialSubtopic bool
	SortDir            string
	PageSize           int
	PageNumber         int
}

// GetFilteredQuestions gets questions with filtering, sorting and pagination
func (s *QuestionService) GetFilteredQuestions(filters QuestionFilters, userID int) ([]models.Question, int, error) {
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
	subjectIDParam := int32(filters.SubjectID)
	if filters.SubjectID <= 0 {
		subjectIDParam = -1 // Special value for "no filter"
	}

	difficultyParam := int32(filters.DifficultyLevel)
	if filters.DifficultyLevel < 0 || filters.DifficultyLevel > 2 {
		difficultyParam = -1 // Special value for "no filter"
	}

	// Log the parameters
	log.Printf("Executing GetFilteredQuestions with filters: %+v, userID: %d", filters, userID)

	// Before executing the query
	var subtopicParam string
	if filters.HasSpecialSubtopic {
		subtopicParam = "Form, Structure, and Sense"
		log.Printf("Using special handling for subtopic: %q", subtopicParam)
	} else {
		subtopicParam = filters.Subtopic
	}

	// Get questions with user data
	dbQuestions, err := s.db.GetFilteredQuestions(ctx, db.GetFilteredQuestionsParams{
		Column1: subjectIDParam,
		Column2: difficultyParam,
		Column3: filters.Topic,
		Column4: subtopicParam,
		Column5: filters.SortDir,
		Limit:   int32(filters.PageSize),
		Offset:  int32(offset),
		UserID:  int32(userID),
	})
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, 0, err
	}

	// Log raw database results in detail
	log.Printf("===== DATABASE QUERY RESULTS =====")
	log.Printf("User ID: %d", userID)
	for i, q := range dbQuestions {
		log.Printf("Question %d (ID: %d):", i, q.ID)
		log.Printf("  - IsSolved: %v (type: %T)", q.IsSolved, q.IsSolved)
		log.Printf("  - IsBookmarked: %v (type: %T)", q.IsBookmarked, q.IsBookmarked)
		log.Printf("  - Incorrect: %v (type: %T)", q.Incorrect, q.Incorrect)
		log.Printf("  - Text: %.30s...", q.QuestionText)
	}
	log.Printf("===== END DATABASE QUERY RESULTS =====")

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
			Passage:            q.Passage.String,
			Bluebook:           q.Bluebook,
			HTMLTable:          q.HtmlTable.String,
			SVGImage:           q.SvgImage.String,
			IsMultipleChoice:   q.IsMultipleChoice.Bool,
			IsSolved:           q.IsSolved,
			IsBookmarked:       q.IsBookmarked,
			Incorrect:          q.Incorrect,
		}
	}

	return questions, totalCount, nil
}
