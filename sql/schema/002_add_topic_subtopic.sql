-- Add topic and subtopic columns to questions table
ALTER TABLE questions 
ADD COLUMN topic VARCHAR(100),
ADD COLUMN subtopic VARCHAR(100); 