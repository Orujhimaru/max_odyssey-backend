package models

import "time"

type Subject struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Question struct {
	ID              int            `json:"id"`
	SubjectID       int            `json:"subject_id"`
	QuestionText    string         `json:"question_text"`
	CorrectAnswer   string         `json:"correct_answer"`
	DifficultyLevel int            `json:"difficulty_level"`
	Explanation     string         `json:"explanation"`
	CreatedAt       time.Time      `json:"created_at"`
	Choices         []AnswerChoice `json:"choices,omitempty"`
	Topic           string         `json:"topic"`
	Subtopic        string         `json:"subtopic"`
	SolveRate       int            `json:"solve_rate"`
}

type AnswerChoice struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	ChoiceText string    `json:"choice_text"`
	IsCorrect  bool      `json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
}
