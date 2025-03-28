-- name: CreateExamResult :one
INSERT INTO exam_results (
    user_id, exam_number, math_score, verbal_score, 
    math_time_taken, verbal_time_taken, exam_data
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetExamResultsByUserID :many
SELECT id, user_id, exam_number, math_score, verbal_score, 
       math_time_taken, verbal_time_taken, exam_data, created_at
FROM exam_results
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetExamResultByID :one
SELECT id, user_id, exam_number, math_score, verbal_score, 
       math_time_taken, verbal_time_taken, exam_data, created_at
FROM exam_results
WHERE id = $1;

-- name: DeleteExamResult :exec
DELETE FROM exam_results
WHERE id = $1 AND user_id = $2; 