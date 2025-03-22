-- Add passage column to questions table
ALTER TABLE questions 
ADD COLUMN passage TEXT DEFAULT '';
