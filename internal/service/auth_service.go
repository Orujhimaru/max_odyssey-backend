package service

import (
	"context"
	"errors"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	db         *db.Queries
	jwtSecret  []byte
	jwtExpires time.Duration
}

func NewAuthService(db *db.Queries, jwtSecret string, jwtExpires time.Duration) *AuthService {
	return &AuthService{
		db:         db,
		jwtSecret:  []byte(jwtSecret),
		jwtExpires: jwtExpires,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username string) (string, error) {
	// Get user from database
	user, err := s.db.GetUserByUsername(context.Background(), username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(s.jwtExpires).Unix(),
	})

	// Sign and return the token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Get the user ID
	userID := int64(claims["user_id"].(float64))

	// Get the user from the database
	user, err := s.db.GetUserByID(context.Background(), int32(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Convert to model
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
