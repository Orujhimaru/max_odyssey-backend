-- name: GetQuestions :many
SELECT id, subject_id, question_text, correct_answer_index, difficulty_level, explanation, 
       created_at, topic, subtopic, solve_rate, choices, passage, bluebook, html_table, svg_image, is_multiple_choice
FROM questions; 

-- name: GetQuestion :one
SELECT id, subject_id, question_text, correct_answer_index, difficulty_level, explanation, 
       created_at, topic, subtopic, solve_rate, choices, passage, bluebook, html_table, svg_image, is_multiple_choice
FROM questions 
WHERE id = $1; 

-- name: GetFilteredQuestions :many
SELECT 
    q.id, 
    q.subject_id, 
    q.question_text, 
    q.correct_answer_index,
    q.difficulty_level, 
    q.explanation, 
    q.created_at, 
    q.topic, 
    q.subtopic, 
    q.solve_rate, 
    q.choices, 
    q.passage, 
    q.bluebook, 
    q.html_table, 
    q.svg_image, 
    q.is_multiple_choice,
    COUNT(*) OVER() AS total_count,
    COALESCE(uq.is_solved, FALSE) as is_solved,
    COALESCE(uq.is_bookmarked, FALSE) as is_bookmarked,
    COALESCE(uq.incorrect, FALSE) as incorrect,
    uq.selected_option
FROM 
    questions q
LEFT JOIN 
    user_questions uq ON q.id = uq.question_id AND uq.user_id = $8
WHERE 
    ($1 = -1 OR q.subject_id = $1) AND
    ($2 = -1 OR q.difficulty_level = $2) AND
    ($3 = '' OR q.topic = $3) AND
    ($4 = '' OR q.subtopic = $4)
ORDER BY 
    CASE WHEN $5 = 'asc' THEN q.solve_rate END ASC,
    CASE WHEN $5 = 'desc' THEN q.solve_rate END DESC
LIMIT $6
OFFSET $7;
