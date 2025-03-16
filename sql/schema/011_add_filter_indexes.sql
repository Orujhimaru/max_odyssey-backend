-- Add indexes for filtering
CREATE INDEX IF NOT EXISTS idx_questions_subject_id ON questions(subject_id);
CREATE INDEX IF NOT EXISTS idx_questions_difficulty_level ON questions(difficulty_level); 