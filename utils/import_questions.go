package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/lib/pq"
)

type JSONQuestion struct {
	ID            int      `json:"id"`
	Type          string   `json:"type"`
	Topic         string   `json:"topic"`
	Subtopic      string   `json:"subtopic"`
	Difficulty    string   `json:"difficulty"`
	Question      string   `json:"question"`
	Passage       string   `json:"passage,omitempty"`
	Choices       []string `json:"choices"`
	CorrectAnswer string   `json:"correctAnswer"`
	Explanation   string   `json:"explanation"`
	SolveRate     int      `json:"solveRate"`
}

// This function reads questions.json, and adds the questions to our database.
func ImportQuestions() {
	// Connect to database
	connStr := "host=localhost port=5431 user=satapp password=satapp123 dbname=sat_tracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read JSON file
	jsonData, err := ioutil.ReadFile("questions.json")
	if err != nil {
		log.Fatal(err)
	}

	var questions []JSONQuestion
	if err := json.Unmarshal(jsonData, &questions); err != nil {
		log.Fatal(err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare statement
	insertQuestion, err := tx.Prepare(`
		INSERT INTO questions (
			subject_id, question_text, difficulty_level, explanation, 
			topic, subtopic, solve_rate, choices, correct_answer_index
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	// Insert questions
	for _, q := range questions {
		// Map difficulty to level (0-2)
		difficultyLevel := map[string]int{
			"Easy":   0,
			"Medium": 1,
			"Hard":   2,
		}[q.Difficulty]

		// Map type to subject_id
		subjectID := map[string]int{
			"Math":   1,
			"Verbal": 2,
		}[q.Type]

		// Combine question and passage if present
		questionText := q.Question
		if q.Passage != "" {
			questionText = fmt.Sprintf("Passage:\n%s\n\nQuestion:\n%s", q.Passage, q.Question)
		}

		// Convert letter answer (A, B, C, D) to index (0, 1, 2, 3)
		correctAnswerIndex := map[string]int{
			"A": 0,
			"B": 1,
			"C": 2,
			"D": 3,
		}[q.CorrectAnswer]

		// Insert question
		var questionID int
		err := insertQuestion.QueryRow(
			subjectID,
			questionText,
			difficultyLevel,
			q.Explanation,
			q.Topic,
			q.Subtopic,
			q.SolveRate,
			pq.Array(q.Choices),
			correctAnswerIndex,
		).Scan(&questionID)
		if err != nil {
			if strings.Contains(err.Error(), "unique_question") {
				// Skip duplicate question
				fmt.Printf("Skipping duplicate question: %s\n", q.Question)
				continue
			}
			tx.Rollback()
			log.Fatal(err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully imported questions!")
}
