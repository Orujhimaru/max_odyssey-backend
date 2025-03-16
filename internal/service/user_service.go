package service

import (
	"context"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/models"
)

type UserService struct {
	db *db.Queries
}

func NewUserService(db *db.Queries) *UserService {
	return &UserService{
		db: db,
	}
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.db.GetUserByID(context.Background(), int32(id))
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:                   int(user.ID),
		Username:             user.Username,
		Role:                 string(user.Role),
		AvatarURL:            user.AvatarUrl.String,
		TargetScore:          int(user.TargetScore.Int32),
		PredictedTotalScore:  int(user.PredictedTotalScore.Int32),
		TotalQuestionsSolved: int(user.TotalQuestionsSolved.Int32),
		CreatedAt:            user.CreatedAt.Time,
	}, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username string, role string) (*models.User, error) {
	// Convert role string to enum
	userRole := db.UserRole(role)

	user, err := s.db.CreateUser(context.Background(), db.CreateUserParams{
		Username: username,
		Role:     userRole,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:                   int(user.ID),
		Username:             user.Username,
		Role:                 string(user.Role),
		AvatarURL:            user.AvatarUrl.String,
		TargetScore:          int(user.TargetScore.Int32),
		PredictedTotalScore:  int(user.PredictedTotalScore.Int32),
		TotalQuestionsSolved: int(user.TotalQuestionsSolved.Int32),
		CreatedAt:            user.CreatedAt.Time,
	}, nil
}
