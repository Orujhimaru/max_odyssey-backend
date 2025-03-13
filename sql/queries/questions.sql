-- name: GetQuestions :many
SELECT id, subject_id, question_text, correct_answer, difficulty_level, explanation, created_at, topic, subtopic, solve_rate 
FROM questions; 