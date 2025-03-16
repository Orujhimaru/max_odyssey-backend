-- First, create the enum type
CREATE TYPE user_role AS ENUM ('free', 'paid', 'admin');

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    role user_role NOT NULL DEFAULT 'free',
    avatar_url TEXT,
    target_score INTEGER,
    predicted_total_score INTEGER,
    total_questions_solved INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

