-- Create user_questions table to track user interactions with questions
CREATE TABLE user_questions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    is_solved BOOLEAN DEFAULT FALSE,
    is_bookmarked BOOLEAN DEFAULT FALSE,
    time_taken INTEGER, -- Time taken in seconds
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure each user can only have one entry per question
    CONSTRAINT unique_user_question UNIQUE (user_id, question_id)
);

-- Create indexes for faster lookups
CREATE INDEX idx_user_questions_user_id ON user_questions(user_id);
CREATE INDEX idx_user_questions_question_id ON user_questions(question_id);
CREATE INDEX idx_user_questions_is_bookmarked ON user_questions(is_bookmarked) WHERE is_bookmarked = TRUE; 