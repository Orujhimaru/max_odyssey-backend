-- Drop the answer_choices table if it exists
DROP TABLE IF EXISTS answer_choices CASCADE;

-- Make sure all columns exist (idempotent)
DO $$
BEGIN
    -- Add topic and subtopic if they don't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                  WHERE table_name='questions' AND column_name='topic') THEN
        ALTER TABLE questions ADD COLUMN topic VARCHAR(100);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                  WHERE table_name='questions' AND column_name='subtopic') THEN
        ALTER TABLE questions ADD COLUMN subtopic VARCHAR(100);
    END IF;
    
    -- Add solve_rate if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                  WHERE table_name='questions' AND column_name='solve_rate') THEN
        ALTER TABLE questions ADD COLUMN solve_rate INTEGER;
    END IF;
    
    -- Add choices array if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                  WHERE table_name='questions' AND column_name='choices') THEN
        ALTER TABLE questions ADD COLUMN choices TEXT[] CHECK (array_length(choices, 1) = 4);
    END IF;
    
    -- Add correct_answer_index if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                  WHERE table_name='questions' AND column_name='correct_answer_index') THEN
        ALTER TABLE questions ADD COLUMN correct_answer_index INTEGER CHECK (correct_answer_index BETWEEN 0 AND 3);
    END IF;
END
$$; 