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

	questions := make([]models.Question, len(dbQuestions))
	for i, q := range dbQuestions {
		questions[i] = models.Question{
			ID:              int(q.ID),
			SubjectID:       int(q.SubjectID.Int32),
			QuestionText:    q.QuestionText,
			Topic:           q.Topic.String,
			Subtopic:        q.Subtopic.String,
			CorrectAnswer:   q.CorrectAnswer,
			DifficultyLevel: int(q.DifficultyLevel.Int32),
			Explanation:     q.Explanation.String,
			CreatedAt:       q.CreatedAt.Time,
			SolveRate:       int(q.SolveRate.Int32),
			Choices:         []models.AnswerChoice{}, // Empty slice for now
		}
	}

	return questions, nil
}

// func (s *QuestionService) CreateQuestion(question db.CreateQuestionParams) error {
// 	return s.db.CreateQuestion(context.Background(), question)
// }
