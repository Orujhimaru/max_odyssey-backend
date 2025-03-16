-- name: GetUserQuestionByIDs :one
SELECT * FROM user_questions
WHERE user_id = $1 AND question_id = $2;

-- name: GetUserBookmarkedQuestions :many
SELECT 
  q.id,
  q.subject_id,
  q.question_text,
  q.difficulty_level,
  q.explanation,
  q.topic,
  q.subtopic,
  q.solve_rate,
  q.choices,
  q.correct_answer_index,
  q.created_at
FROM questions q
JOIN user_questions uq ON q.id = uq.question_id
WHERE uq.user_id = $1 AND uq.is_bookmarked = TRUE;

-- name: GetUserSolvedQuestions :many
SELECT 
  q.id,
  q.subject_id,
  q.question_text,
  q.difficulty_level,
  q.explanation,
  q.topic,
  q.subtopic,
  q.solve_rate,
  q.choices,
  q.correct_answer_index,
  q.created_at
FROM questions q
JOIN user_questions uq ON q.id = uq.question_id
WHERE uq.user_id = $1 AND uq.is_solved = TRUE;

-- name: CreateUserQuestion :one
INSERT INTO user_questions (
    user_id, question_id, is_solved, is_bookmarked, time_taken
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (user_id, question_id) 
DO UPDATE SET
    is_solved = EXCLUDED.is_solved,
    is_bookmarked = EXCLUDED.is_bookmarked,
    time_taken = EXCLUDED.time_taken
RETURNING *;

-- name: ToggleBookmark :one
UPDATE user_questions
SET is_bookmarked = NOT is_bookmarked
WHERE user_id = $1 AND question_id = $2
RETURNING *;

-- name: MarkQuestionSolved :one
UPDATE user_questions
SET 
    is_solved = TRUE,
    time_taken = $3
WHERE user_id = $1 AND question_id = $2
RETURNING *; 