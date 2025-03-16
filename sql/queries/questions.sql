-- name: GetQuestions :many
SELECT id, subject_id, question_text, correct_answer_index, difficulty_level, explanation, 
       created_at, topic, subtopic, solve_rate, choices
FROM questions; 

-- name: GetQuestion :one
SELECT id, subject_id, question_text, correct_answer_index, difficulty_level, explanation, 
       created_at, topic, subtopic, solve_rate, choices
FROM questions 
WHERE id = $1; 