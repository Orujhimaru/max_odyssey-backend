-- Add is_multiple_choice column to questions table
ALTER TABLE questions ADD COLUMN is_multiple_choice BOOLEAN DEFAULT TRUE;

-- Update existing records to have is_multiple_choice = true
UPDATE questions SET is_multiple_choice = TRUE; 