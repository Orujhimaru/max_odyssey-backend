-- Add index on solve_rate for faster sorting and filtering
CREATE INDEX idx_questions_solve_rate ON questions(solve_rate); 