-- Drop the existing unique constraint
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'unique_question' 
        AND conrelid = 'questions'::regclass::oid
    ) THEN
        ALTER TABLE questions DROP CONSTRAINT unique_question;
    END IF;
END $$;

-- Add a new unique constraint on question_text and explanation
ALTER TABLE questions ADD CONSTRAINT questions_text_explanation_unique 
UNIQUE (question_text, explanation);

-- This allows questions with the same text but different explanations to be inserted
-- Or questions with the same explanation but different text 