package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"max-odyssey-backend/internal/db"

	"github.com/sqlc-dev/pqtype"
)

type ExamResultService struct {
	queries *db.Queries
}

func NewExamResultService(queries *db.Queries) *ExamResultService {
	return &ExamResultService{
		queries: queries,
	}
}

type ExamResultResponse struct {
	ID              int64           `json:"id"`
	UserID          int32           `json:"user_id"`
	ExamNumber      int32           `json:"exam_number"`
	MathScore       sql.NullInt32   `json:"math_score"`
	VerbalScore     sql.NullInt32   `json:"verbal_score"`
	MathTimeTaken   sql.NullInt32   `json:"math_time_taken"`
	VerbalTimeTaken sql.NullInt32   `json:"verbal_time_taken"`
	ExamData        json.RawMessage `json:"exam_data"`
	CreatedAt       sql.NullTime    `json:"created_at"`
}

type CreateExamResultRequest struct {
	UserID          int64           `json:"user_id"`
	ExamNumber      int32           `json:"exam_number"`
	MathScore       *int32          `json:"math_score"`
	VerbalScore     *int32          `json:"verbal_score"`
	MathTimeTaken   *int32          `json:"math_time_taken"`
	VerbalTimeTaken *int32          `json:"verbal_time_taken"`
	ExamData        json.RawMessage `json:"exam_data"`
}

func (s *ExamResultService) CreateExamResult(ctx context.Context, req CreateExamResultRequest) (*ExamResultResponse, error) {
	// Convert nullable fields
	mathScore := sql.NullInt32{}
	if req.MathScore != nil {
		mathScore.Int32 = *req.MathScore
		mathScore.Valid = true
	}

	verbalScore := sql.NullInt32{}
	if req.VerbalScore != nil {
		verbalScore.Int32 = *req.VerbalScore
		verbalScore.Valid = true
	}

	mathTimeTaken := sql.NullInt32{}
	if req.MathTimeTaken != nil {
		mathTimeTaken.Int32 = *req.MathTimeTaken
		mathTimeTaken.Valid = true
	}

	verbalTimeTaken := sql.NullInt32{}
	if req.VerbalTimeTaken != nil {
		verbalTimeTaken.Int32 = *req.VerbalTimeTaken
		verbalTimeTaken.Valid = true
	}

	// Create exam result in database
	result, err := s.queries.CreateExamResult(ctx, db.CreateExamResultParams{
		UserID:          int32(req.UserID),
		ExamNumber:      req.ExamNumber,
		MathScore:       mathScore,
		VerbalScore:     verbalScore,
		MathTimeTaken:   mathTimeTaken,
		VerbalTimeTaken: verbalTimeTaken,
		ExamData:        pqtype.NullRawMessage{RawMessage: req.ExamData, Valid: len(req.ExamData) > 0},
	})
	if err != nil {
		return nil, err
	}

	return &ExamResultResponse{
		ID:              int64(result.ID),
		UserID:          result.UserID,
		ExamNumber:      result.ExamNumber,
		MathScore:       result.MathScore,
		VerbalScore:     result.VerbalScore,
		MathTimeTaken:   result.MathTimeTaken,
		VerbalTimeTaken: result.VerbalTimeTaken,
		ExamData:        result.ExamData.RawMessage,
		CreatedAt:       result.CreatedAt,
	}, nil
}

func (s *ExamResultService) GetExamResultsByUserID(ctx context.Context, userID int32) ([]ExamResultResponse, error) {
	results, err := s.queries.GetExamResultsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []ExamResultResponse
	for _, result := range results {
		response = append(response, ExamResultResponse{
			ID:              int64(result.ID),
			UserID:          result.UserID,
			ExamNumber:      result.ExamNumber,
			MathScore:       result.MathScore,
			VerbalScore:     result.VerbalScore,
			MathTimeTaken:   result.MathTimeTaken,
			VerbalTimeTaken: result.VerbalTimeTaken,
			ExamData:        result.ExamData.RawMessage,
			CreatedAt:       result.CreatedAt,
		})
	}

	return response, nil
}

func (s *ExamResultService) GetExamResultByID(ctx context.Context, id int64) (*ExamResultResponse, error) {
	result, err := s.queries.GetExamResultByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &ExamResultResponse{
		ID:              int64(result.ID),
		UserID:          result.UserID,
		ExamNumber:      result.ExamNumber,
		MathScore:       result.MathScore,
		VerbalScore:     result.VerbalScore,
		MathTimeTaken:   result.MathTimeTaken,
		VerbalTimeTaken: result.VerbalTimeTaken,
		ExamData:        result.ExamData.RawMessage,
		CreatedAt:       result.CreatedAt,
	}, nil
}

func (s *ExamResultService) DeleteExamResult(ctx context.Context, id int64, userID int32) error {
	return s.queries.DeleteExamResult(ctx, db.DeleteExamResultParams{
		ID:     int32(id),
		UserID: userID,
	})
}
