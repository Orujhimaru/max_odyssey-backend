-- Add default value of -1 to selected_option column
ALTER TABLE user_questions 
ALTER COLUMN selected_option SET DEFAULT -1;

-- Update existing NULL values to -1
UPDATE user_questions 
SET selected_option = -1 
WHERE selected_option IS NULL; 