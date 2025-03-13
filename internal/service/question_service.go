package service

import (
	"context"
	"max-odyssey-backend/internal/db"
)

type QuestionService struct {
	queries *db.Queries
}

func NewQuestionService(queries *db.Queries) *QuestionService {
	return &QuestionService{queries: queries}
}

func (s *QuestionService) GetQuestions(ctx context.Context) ([]db.Question, error) {
	return s.queries.GetQuestions(ctx)
}
