-- name: GetQuestions :many
SELECT id, subject_id, question_text, correct_answer_index, difficulty_level, explanation, 
       created_at, topic, subtopic, solve_rate, choices
FROM questions; 

-- name: GetQuestion :one
SELECT id, subject_id, question_text, correct_answer_index, difficulty_level, explanation, 
       created_at, topic, subtopic, solve_rate, choices
FROM questions 
WHERE id = $1; 

-- name: GetFilteredQuestions :many
SELECT 
  id, subject_id, question_text, correct_answer_index, 
  difficulty_level, explanation, topic, subtopic, solve_rate, choices,
  COUNT(*) OVER() AS total_count
FROM questions
WHERE 
  ($1 = -1 OR subject_id = $1) AND 
  ($2 = -1 OR difficulty_level = $2) AND
  ($3 = '' OR topic = $3) AND
  ($4 = '' OR subtopic = $4)
ORDER BY 
  CASE WHEN $5 = 'asc' THEN solve_rate END ASC,
  CASE WHEN $5 = 'desc' THEN solve_rate END DESC
LIMIT $6 OFFSET $7;