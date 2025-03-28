 CREATE TABLE IF NOT EXISTS exam_results (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    exam_number INT NOT NULL,
    math_score INT,
    verbal_score INT,
    math_time_taken INT,
    verbal_time_taken INT,
    exam_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_exam_results_user_id ON exam_results(user_id);