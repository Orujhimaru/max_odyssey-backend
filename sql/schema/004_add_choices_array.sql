
-- Add choices array and correct_answer_index to questions
ALTER TABLE questions 
ADD COLUMN choices TEXT[] CHECK (array_length(choices, 1) = 4),
ADD COLUMN correct_answer_index INTEGER CHECK (correct_answer_index BETWEEN 0 AND 3); 