-- Add bluebook boolean field to questions table
ALTER TABLE questions 
ADD COLUMN bluebook BOOLEAN DEFAULT FALSE; 