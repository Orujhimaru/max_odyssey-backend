package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
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
}

func main() {
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

	// Prepare statements
	insertQuestion, err := tx.Prepare(`
		INSERT INTO questions (subject_id, question_text, correct_answer, difficulty_level, explanation)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	insertChoice, err := tx.Prepare(`
		INSERT INTO answer_choices (question_id, choice_text, is_correct)
		VALUES ($1, $2, $3)`)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	// Insert questions
	for _, q := range questions {
		// Map difficulty to level (1-5)
		difficultyLevel := map[string]int{
			"Easy":   1,
			"Medium": 3,
			"Hard":   5,
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

		// Insert question
		var questionID int
		err := insertQuestion.QueryRow(
			subjectID,
			questionText,
			q.CorrectAnswer,
			difficultyLevel,
			q.Explanation,
		).Scan(&questionID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		// Insert choices
		for _, choice := range q.Choices {
			isCorrect := strings.HasPrefix(choice, q.CorrectAnswer+")")
			_, err = insertChoice.Exec(questionID, choice, isCorrect)
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully imported questions!")
}
