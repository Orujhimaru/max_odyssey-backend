-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (
    username, role, avatar_url, target_score
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUserScore :one
UPDATE users
SET 
    predicted_total_score = $2,
    total_questions_solved = total_questions_solved + 1
WHERE id = $1
RETURNING *; 