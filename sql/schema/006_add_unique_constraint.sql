-- Add a unique constraint on question_text and topic
ALTER TABLE questions 
ADD CONSTRAINT unique_question UNIQUE (question_text, topic, subtopic); 