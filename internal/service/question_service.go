package service

import (
	"context"
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
