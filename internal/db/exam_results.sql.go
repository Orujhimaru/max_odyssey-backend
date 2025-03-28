// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: exam_results.sql

package db

import (
	"context"
	"database/sql"
	"github.com/sqlc-dev/pqtype"
	"log"
)

const createExamResult = `-- name: CreateExamResult :one
INSERT INTO exam_results (
    user_id, exam_number, math_score, verbal_score, 
    math_time_taken, verbal_time_taken, exam_data
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, user_id, exam_number, math_score, verbal_score, math_time_taken, verbal_time_taken, exam_data, created_at
`

type CreateExamResultParams struct {
	UserID          int32
	ExamNumber      int32
	MathScore       sql.NullInt32
	VerbalScore     sql.NullInt32
	MathTimeTaken   sql.NullInt32
	VerbalTimeTaken sql.NullInt32
	ExamData        pqtype.NullRawMessage
}

func (q *Queries) CreateExamResult(ctx context.Context, arg CreateExamResultParams) (ExamResult, error) {
	row := q.db.QueryRowContext(ctx, createExamResult,
		arg.UserID,
		arg.ExamNumber,
		arg.MathScore,
		arg.VerbalScore,
		arg.MathTimeTaken,
		arg.VerbalTimeTaken,
		arg.ExamData,
	)
	var i ExamResult
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExamNumber,
		&i.MathScore,
		&i.VerbalScore,
		&i.MathTimeTaken,
		&i.VerbalTimeTaken,
		&i.ExamData,
		&i.CreatedAt,
	)
	return i, err
}

const deleteExamResult = `-- name: DeleteExamResult :exec
DELETE FROM exam_results
WHERE id = $1 AND user_id = $2
`

type DeleteExamResultParams struct {
	ID     int32
	UserID int32
}

func (q *Queries) DeleteExamResult(ctx context.Context, arg DeleteExamResultParams) error {
	_, err := q.db.ExecContext(ctx, deleteExamResult, arg.ID, arg.UserID)
	return err
}

const getExamResultByID = `-- name: GetExamResultByID :one
SELECT id, user_id, exam_number, math_score, verbal_score, 
       math_time_taken, verbal_time_taken, exam_data, created_at
FROM exam_results
WHERE id = $1
`

func (q *Queries) GetExamResultByID(ctx context.Context, id int32) (ExamResult, error) {
	row := q.db.QueryRowContext(ctx, getExamResultByID, id)
	var i ExamResult
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExamNumber,
		&i.MathScore,
		&i.VerbalScore,
		&i.MathTimeTaken,
		&i.VerbalTimeTaken,
		&i.ExamData,
		&i.CreatedAt,
	)
	return i, err
}

const getExamResultsByUserID = `-- name: GetExamResultsByUserID :many
SELECT id, user_id, exam_number, math_score, verbal_score, 
       math_time_taken, verbal_time_taken, exam_data, created_at
FROM exam_results
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetExamResultsByUserID(ctx context.Context, userID int32) ([]ExamResult, error) {
	log.Printf("SQL Query: GetExamResultsByUserID with userID=%d", userID)
	rows, err := q.db.QueryContext(ctx, getExamResultsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExamResult
	for rows.Next() {
		var i ExamResult
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ExamNumber,
			&i.MathScore,
			&i.VerbalScore,
			&i.MathTimeTaken,
			&i.VerbalTimeTaken,
			&i.ExamData,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
