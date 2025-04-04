package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Question represents a question in the system
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
	Passage            string    `json:"passage"`
	Bluebook           bool      `json:"bluebook"`
	HTMLTable          string    `json:"html_table"`
	SVGImage           string    `json:"svg_image"`
	IsMultipleChoice   bool      `json:"is_multiple_choice"`
	IsSolved           bool      `json:"is_solved"`
	IsBookmarked       bool      `json:"is_bookmarked"`
	Incorrect          bool      `json:"incorrect"`
	SelectedOption     *int      `json:"selected_option,omitempty"`
}

// User represents a user in the system
type User struct {
	ID                   int       `json:"id"`
	Username             string    `json:"username"`
	Role                 string    `json:"role"`
	AvatarURL            string    `json:"avatar_url"`
	TargetScore          int       `json:"target_score"`
	PredictedTotalScore  int       `json:"predicted_total_score"`
	TotalQuestionsSolved int       `json:"total_questions_solved"`
	CreatedAt            time.Time `json:"created_at"`
}

// UserQuestion represents a user's interaction with a question
type UserQuestion struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	QuestionID     int       `json:"question_id"`
	IsSolved       bool      `json:"is_solved"`
	IsBookmarked   bool      `json:"is_bookmarked"`
	TimeTaken      int       `json:"time_taken"`
	CreatedAt      time.Time `json:"created_at"`
	Incorrect      bool      `json:"incorrect"`
	SelectedOption *int      `json:"selected_option,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserSkill represents a user's skill level
type UserSkill struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	SkillName  string    `json:"skill_name"`
	SkillScore float64   `json:"skill_score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ExamResult represents a user's exam result
type ExamResult struct {
	ID              int64           `json:"id"`
	UserID          int32           `json:"user_id"`
	ExamNumber      int32           `json:"exam_number"`
	MathScore       sql.NullInt32   `json:"math_score"`
	VerbalScore     sql.NullInt32   `json:"verbal_score"`
	MathTimeTaken   sql.NullInt32   `json:"math_time_taken"`
	VerbalTimeTaken sql.NullInt32   `json:"verbal_time_taken"`
	ExamData        json.RawMessage `json:"exam_data"`
	CreatedAt       time.Time       `json:"created_at"`
}

// QuestionUpdateRequest represents a request to update a user's interaction with a question
type QuestionUpdateItem struct {
	QuestionID     int  `json:"questionId"`
	IsBookmarked   bool `json:"isBookmarked"`
	IsSolved       bool `json:"isSolved"`
	IsIncorrect    bool `json:"isIncorrect"`
	SelectedOption *int `json:"selectedOption"`
}

// BatchQuestionUpdateRequest represents a batch update request for multiple questions
type BatchQuestionUpdateRequest struct {
	Questions []QuestionUpdateItem `json:"questions"`
}
