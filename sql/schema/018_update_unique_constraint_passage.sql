-- Drop the existing unique constraint
ALTER TABLE questions 
DROP CONSTRAINT IF EXISTS unique_question;

-- Add a new unique constraint using passage instead of question_text
ALTER TABLE questions 
ADD CONSTRAINT unique_question UNIQUE (question_text, passage, topic, subtopic); 