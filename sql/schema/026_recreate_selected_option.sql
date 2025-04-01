-- Drop the column if it exists
ALTER TABLE user_questions DROP COLUMN IF EXISTS selected_option;

-- Add the column fresh
ALTER TABLE user_questions ADD COLUMN selected_option INT DEFAULT -1; 