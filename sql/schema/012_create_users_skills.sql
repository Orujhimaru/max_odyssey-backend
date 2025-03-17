-- Create users_skills table to track user skill levels
CREATE TABLE users_skills (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_name VARCHAR(100) NOT NULL,
    skill_score REAL NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure each user can only have one entry per skill
    CONSTRAINT unique_user_skill UNIQUE (user_id, skill_name)
);

-- Create indexes for faster lookups
CREATE INDEX idx_users_skills_user_id ON users_skills(user_id);
CREATE INDEX idx_users_skills_skill_name ON users_skills(skill_name); 