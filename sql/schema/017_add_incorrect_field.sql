-- Add incorrect boolean field to user_questions table
ALTER TABLE user_questions 
ADD COLUMN incorrect BOOLEAN DEFAULT FALSE NOT NULL; 