CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    subject_id INTEGER REFERENCES subjects(id),
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    difficulty_level INTEGER CHECK (difficulty_level >= 0 AND difficulty_level <= 2),
    explanation TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE answer_choices (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id),
    choice_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add some initial subjects
INSERT INTO subjects (name) VALUES 
    ('Math - No Calculator'),
    ('Math - Calculator'),
    ('Reading'),
    ('Writing and Language'); 