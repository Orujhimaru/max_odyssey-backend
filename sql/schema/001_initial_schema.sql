CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    subject_id INTEGER REFERENCES subjects(id),
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    difficulty_level INTEGER CHECK (difficulty_level >= 0 AND difficulty_level <= 2),
    explanation TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

