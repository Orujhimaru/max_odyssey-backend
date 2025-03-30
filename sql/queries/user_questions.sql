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
  q.created_at,
  q.passage,
  COUNT(*) OVER() AS total_count
FROM questions q
JOIN user_questions uq ON q.id = uq.question_id
WHERE uq.user_id = $1 AND uq.is_bookmarked = TRUE
ORDER BY 
  CASE 
    WHEN $2 = 'desc' THEN q.solve_rate * -1  -- Multiply by -1 for descending
    ELSE q.solve_rate
  END;

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

-- name: GetUserQuestions :many
SELECT id, user_id, question_id, is_solved, is_bookmarked, time_taken, created_at, incorrect
FROM user_questions
WHERE user_id = $1;

-- name: GetUserQuestion :one
SELECT id, user_id, question_id, is_solved, is_bookmarked, time_taken, created_at, incorrect
FROM user_questions
WHERE user_id = $1 AND question_id = $2;

-- name: CreateUserQuestion :one
INSERT INTO user_questions (
  user_id, question_id, is_solved, is_bookmarked, time_taken, incorrect
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, user_id, question_id, is_solved, is_bookmarked, time_taken, created_at, incorrect;

-- name: ToggleBookmark :one
UPDATE user_questions
SET is_bookmarked = NOT is_bookmarked
WHERE user_id = $1 AND question_id = $2
RETURNING *;

-- name: MarkQuestionSolved :one
UPDATE user_questions
SET 
    is_solved = TRUE,
    time_taken = $3,
    incorrect = $4
WHERE user_id = $1 AND question_id = $2
RETURNING *;

-- name: ToggleSolved :one
UPDATE user_questions
SET is_solved = NOT is_solved
WHERE user_id = $1 AND question_id = $2
RETURNING *;

-- name: GetUserBookmarkedQuestionsAsc :many
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
  q.created_at,
  q.passage,
  COUNT(*) OVER() AS total_count
FROM questions q
JOIN user_questions uq ON q.id = uq.question_id
WHERE uq.user_id = $1 AND uq.is_bookmarked = TRUE
ORDER BY q.solve_rate ASC;

-- name: GetUserBookmarkedQuestionsDesc :many
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
  q.created_at,
  q.passage,
  COUNT(*) OVER() AS total_count
FROM questions q
JOIN user_questions uq ON q.id = uq.question_id
WHERE uq.user_id = $1 AND uq.is_bookmarked = TRUE
ORDER BY q.solve_rate DESC;

-- name: UpdateUserQuestionData :one
UPDATE user_questions
SET 
    is_solved = $3,
    is_bookmarked = $4,
    time_taken = $5,
    incorrect = $6
WHERE user_id = $1 AND question_id = $2
RETURNING *;

-- name: CheckUserQuestionExists :one
SELECT EXISTS(
  SELECT 1 FROM user_questions 
  WHERE user_id = $1 AND question_id = $2
) AS exists;

-- name: UpdateUserQuestion :exec
UPDATE user_questions
SET is_solved = $3, is_bookmarked = $4, incorrect = $5
WHERE user_id = $1 AND question_id = $2;

