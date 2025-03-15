package models

import "time"


type Question struct {
	ID                 int       `json:"id"`
	SubjectID          int       `json:"subject_id"`
	QuestionText       string    `json:"question_text"`
	Choices            []string  `json:"choices"`
	CorrectAnswerIndex int       `json:"correct_answer_index"`
	DifficultyLevel    int       `json:"difficulty_level"`
	Explanation        string    `json:"explanation"`
	Topic              string    `json:"topic"`
	Subtopic           string    `json:"subtopic"`
	SolveRate          int       `json:"solve_rate"`
	CreatedAt          time.Time `json:"created_at"`
}


