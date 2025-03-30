package service

import (
	"context"
	"database/sql"
	"log"
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
func (s *UserQuestionService) GetBookmarkedQuestions(userID int64, sortDir string) ([]models.Question, int, error) {
	// Set default sort direction if not provided
	if sortDir == "" {
		sortDir = "asc"
	}

	// Get bookmarked questions from database
	var totalCount int
	var questions []models.Question

	// Call the appropriate function based on sort direction
	if sortDir == "desc" {
		log.Printf("QUERY PARAMS: GetUserBookmarkedQuestionsDesc - userID=%d", userID)
		dbQuestions, err := s.db.GetUserBookmarkedQuestionsDesc(context.Background(), int32(userID))
		if err != nil {
			log.Printf("QUERY ERROR: %v", err)
			return nil, 0, err
		}

		// Extract total count from first row
		if len(dbQuestions) > 0 {
			totalCount = int(dbQuestions[0].TotalCount)
		} else {
			totalCount = 0
		}

		// Convert to models
		questions = make([]models.Question, len(dbQuestions))
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
			}
		}
	} else {
		log.Printf("QUERY PARAMS: GetUserBookmarkedQuestionsAsc - userID=%d", userID)
		dbQuestions, err := s.db.GetUserBookmarkedQuestionsAsc(context.Background(), int32(userID))
		if err != nil {
			log.Printf("QUERY ERROR: %v", err)
			return nil, 0, err
		}

		// Extract total count from first row
		if len(dbQuestions) > 0 {
			totalCount = int(dbQuestions[0].TotalCount)
		} else {
			totalCount = 0
		}

		// Convert to models
		questions = make([]models.Question, len(dbQuestions))
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
			}
		}
	}

	return questions, totalCount, nil
}

// ToggleBookmark toggles the bookmark status of a question for a user
func (s *UserQuestionService) ToggleBookmark(userID, questionID int64) error {
	// First, check if the question exists
	log.Printf("Checking if question %d exists", questionID)
	_, err := s.db.GetQuestion(context.Background(), int32(questionID))
	if err != nil {
		log.Printf("Error checking question: %v", err)
		return err
	}

	// Now check if a user_question record exists
	log.Printf("Checking if user_question record exists for user %d and question %d", userID, questionID)
	_, err = s.db.GetUserQuestionByIDs(context.Background(), db.GetUserQuestionByIDsParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
	})

	if err != nil {
		if err == sql.ErrNoRows {
			// Record doesn't exist, create a new one with is_bookmarked=true
			log.Printf("Creating new user_question record with bookmark=true")
			_, err = s.db.CreateUserQuestion(context.Background(), db.CreateUserQuestionParams{
				UserID:       int32(userID),
				QuestionID:   int32(questionID),
				IsSolved:     sql.NullBool{Bool: false, Valid: true},
				IsBookmarked: sql.NullBool{Bool: true, Valid: true},
				TimeTaken:    sql.NullInt32{Valid: false},
				Incorrect:    false,
			})
			return err
		}
		log.Printf("Error checking user_question: %v", err)
		return err
	}

	// Record exists, toggle the bookmark
	log.Printf("Toggling bookmark for existing record")
	_, err = s.db.ToggleBookmark(context.Background(), db.ToggleBookmarkParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
	})
	if err != nil {
		log.Printf("Error toggling bookmark: %v", err)
	}
	return err
}

// MarkQuestionSolved marks a question as solved for a user
func (s *UserQuestionService) MarkQuestionSolved(userID, questionID int64, timeTaken int32, incorrect bool) error {
	_, err := s.db.MarkQuestionSolved(context.Background(), db.MarkQuestionSolvedParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
		TimeTaken:  sql.NullInt32{Int32: timeTaken, Valid: true},
		Incorrect:  incorrect,
	})
	return err
}

// GetQuestionsByDifficulty gets questions sorted by solve rate (ascending or descending)
func (s *UserQuestionService) GetQuestionsByDifficulty(ascending bool) ([]models.Question, error) {
	// You'll need to add a new query for this in your SQL files
	// For now, we'll just return an empty slice
	return []models.Question{}, nil
}

// ToggleSolved toggles the solved status of a question for a user
func (s *UserQuestionService) ToggleSolved(userID, questionID int64) error {
	// First, check if the question exists
	log.Printf("Checking if question %d exists", questionID)
	_, err := s.db.GetQuestion(context.Background(), int32(questionID))
	if err != nil {
		log.Printf("Error checking question: %v", err)
		return err
	}

	// Now check if a user_question record exists
	log.Printf("Checking if user_question record exists for user %d and question %d", userID, questionID)
	_, err = s.db.GetUserQuestionByIDs(context.Background(), db.GetUserQuestionByIDsParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
	})

	if err != nil {
		if err == sql.ErrNoRows {
			// Record doesn't exist, create a new one with is_solved=true
			log.Printf("Creating new user_question record with solved=true")
			_, err = s.db.CreateUserQuestion(context.Background(), db.CreateUserQuestionParams{
				UserID:       int32(userID),
				QuestionID:   int32(questionID),
				IsSolved:     sql.NullBool{Bool: true, Valid: true},
				IsBookmarked: sql.NullBool{Bool: false, Valid: true},
				TimeTaken:    sql.NullInt32{Valid: false},
				Incorrect:    false,
			})
			return err
		}
		log.Printf("Error checking user_question: %v", err)
		return err
	}

	// Record exists, toggle the solved status
	log.Printf("Toggling solved status for existing record")
	_, err = s.db.ToggleSolved(context.Background(), db.ToggleSolvedParams{
		UserID:     int32(userID),
		QuestionID: int32(questionID),
	})
	if err != nil {
		log.Printf("Error toggling solved status: %v", err)
	}
	return err
}

// UpdateUserQuestionData updates all fields for a user question
func (s *UserQuestionService) UpdateUserQuestionData(userID, questionID int64, isSolved bool, isBookmarked bool, timeTaken int32, incorrect bool) error {
	_, err := s.db.UpdateUserQuestionData(context.Background(), db.UpdateUserQuestionDataParams{
		UserID:       int32(userID),
		QuestionID:   int32(questionID),
		IsSolved:     sql.NullBool{Bool: isSolved, Valid: true},
		IsBookmarked: sql.NullBool{Bool: isBookmarked, Valid: true},
		TimeTaken:    sql.NullInt32{Int32: timeTaken, Valid: true},
		Incorrect:    incorrect,
	})
	return err
}

// BatchUpdateUserQuestions updates multiple user questions at once
func (s *UserQuestionService) BatchUpdateUserQuestions(ctx context.Context, userID int32, updates []models.QuestionUpdateItem) error {
	log.Printf("Starting batch update for user %d with %d questions", userID, len(updates))

	// Process each update individually without a transaction
	for i, update := range updates {
		log.Printf("Processing update %d/%d: Question ID %d", i+1, len(updates), update.QuestionID)

		// Check if the user question already exists
		exists, err := s.db.CheckUserQuestionExists(ctx, db.CheckUserQuestionExistsParams{
			UserID:     userID,
			QuestionID: int32(update.QuestionID),
		})
		if err != nil {
			log.Printf("Error checking if question exists: %v", err)
			return err
		}

		if exists {
			// Update existing record
			log.Printf("Question %d exists for user %d, updating record", update.QuestionID, userID)
			// First get the current user question to preserve the bookmark status
			userQuestion, err := s.db.GetUserQuestion(ctx, db.GetUserQuestionParams{
				UserID:     userID,
				QuestionID: int32(update.QuestionID),
			})
			if err != nil {
				log.Printf("Error getting current user question: %v", err)
				return err
			}

			err = s.db.UpdateUserQuestion(ctx, db.UpdateUserQuestionParams{
				UserID:       userID,
				QuestionID:   int32(update.QuestionID),
				IsSolved:     sql.NullBool{Bool: update.IsSolved, Valid: true},
				IsBookmarked: userQuestion.IsBookmarked,
				Incorrect:    update.IsIncorrect,
			})
			if err != nil {
				log.Printf("Error updating question %d: %v", update.QuestionID, err)
				return err
			}
			log.Printf("Successfully updated question %d for user %d", update.QuestionID, userID)
		} else {
			// Create new record
			log.Printf("Question %d does not exist for user %d, creating new record", update.QuestionID, userID)
			_, err = s.db.CreateUserQuestion(ctx, db.CreateUserQuestionParams{
				UserID:       userID,
				QuestionID:   int32(update.QuestionID),
				IsSolved:     sql.NullBool{Bool: update.IsSolved, Valid: true},
				TimeTaken:    sql.NullInt32{Valid: false},
				Incorrect:    update.IsIncorrect,
				IsBookmarked: sql.NullBool{Bool: false, Valid: true},
			})
			if err != nil {
				log.Printf("Error creating record for question %d: %v", update.QuestionID, err)
				return err
			}
			log.Printf("Successfully created record for question %d for user %d", update.QuestionID, userID)
		}
	}

	log.Printf("Batch update completed successfully for user %d, updated %d questions", userID, len(updates))
	return nil
}
