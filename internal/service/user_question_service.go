package service

import (
	"context"
	"database/sql"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/models"
)

type UserQuestionService struct {
	db *db.Queries
}

func NewUserQuestionService(db *db.Queries) *UserQuestionService {
	return &UserQuestionService{
		db: db,
	}
}

// GetBookmarkedQuestions gets all bookmarked questions for a user
func (s *UserQuestionService) GetBookmarkedQuestions(userID int64) ([]models.Question, error) {
	dbQuestions, err := s.db.GetUserBookmarkedQuestions(context.Background(), int32(userID))
	if err != nil {
		return nil, err
	}

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
			Choices:            []string{}, // You'll need to convert the array
		}
	}

	return questions, nil
}

// ToggleBookmark toggles the bookmark status of a question for a user
func (s *UserQuestionService) ToggleBookmark(userID, questionID int64) error {
	_, err := s.db.ToggleBookmark(context.Background(), db.ToggleBookmarkParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
	})
	return err
}

// MarkQuestionSolved marks a question as solved for a user
func (s *UserQuestionService) MarkQuestionSolved(userID, questionID int64, timeTaken int32) error {
	_, err := s.db.MarkQuestionSolved(context.Background(), db.MarkQuestionSolvedParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
		TimeTaken:  sql.NullInt32{Int32: timeTaken, Valid: true},
	})
	return err
}

// GetQuestionsByDifficulty gets questions sorted by solve rate (ascending or descending)
func (s *UserQuestionService) GetQuestionsByDifficulty(ascending bool) ([]models.Question, error) {
	// You'll need to add a new query for this in your SQL files
	// For now, we'll just return an empty slice
	return []models.Question{}, nil
}
