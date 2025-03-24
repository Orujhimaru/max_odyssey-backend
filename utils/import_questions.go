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

// JSONQuestion represents a question in the JSON file
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
	Bluebook      bool     `json:"bluebook"`
	HTMLTable     string   `json:"html_table,omitempty"`
	SVGImage      string   `json:"svg_image,omitempty"`
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

	// Add this before the transaction begins
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Database has %d existing questions\n", count)

	// Read JSON file
	jsonData, err := ioutil.ReadFile("questions.json")
	if err != nil {
		log.Fatal(err)
	}

	var questions []JSONQuestion
	if err := json.Unmarshal(jsonData, &questions); err != nil {
		log.Fatal(err)
	}

	// // Check for duplicates in the JSON file
	// uniqueCheck := make(map[string]bool)
	// duplicatesInJSON := 0

	// for i, q := range questions {
	// 	// Create a unique key based on the constraint fields
	// 	key := fmt.Sprintf("%s|%s|%s", q.QuestionText, q.Topic, q.Subtopic)

	// 	if uniqueCheck[key] {
	// 		// Truncate question text for display if needed
	// 		displayText := q.QuestionText
	// 		if len(displayText) > 50 {
	// 			displayText = displayText[:50] + "..."
	// 		}

	// 		fmt.Printf("Duplicate found in JSON (index #%d): %s\n", i+1, displayText)
	// 		fmt.Printf("  Topic: %s, Subtopic: %s\n", q.Topic, q.Subtopic)
	// 		duplicatesInJSON++
	// 	}
	// 	uniqueCheck[key] = true
	// }

	// fmt.Printf("Found %d duplicates in JSON file\n", duplicatesInJSON)

	// Insert questions
	for i, q := range questions {
		// Start a new transaction for each question
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		// Prepare statement
		insertQuestion, err := tx.Prepare(`
			INSERT INTO questions (
				subject_id, question_text, difficulty_level, explanation, 
				topic, subtopic, solve_rate, choices, correct_answer_index, passage, bluebook, html_table, svg_image
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING id`)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

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

		// Get question text
		questionText := q.Question

		// Convert letter answer (A, B, C, D) to index (0, 1, 2, 3)
		correctAnswerIndex := map[string]int{
			"A": 0,
			"B": 1,
			"C": 2,
			"D": 3,
		}[q.CorrectAnswer]

		// Insert question
		var questionID int
		err = insertQuestion.QueryRow(
			subjectID,
			questionText,
			difficultyLevel,
			q.Explanation,
			q.Topic,
			q.Subtopic,
			q.SolveRate,
			pq.Array(q.Choices),
			correctAnswerIndex,
			q.Passage,
			q.Bluebook,
			q.HTMLTable,
			q.SVGImage,
		).Scan(&questionID)

		if err != nil {
			if strings.Contains(err.Error(), "unique_question") {
				// Log duplicate
				displayText := "Empty question text"
				if len(questionText) > 0 {
					if len(questionText) > 50 {
						displayText = questionText[:50] + "..."
					} else {
						displayText = questionText
					}
				}

				fmt.Printf("Skipping duplicate question #%d: %s\n", i+1, displayText)
				tx.Rollback() // Roll back this transaction
				continue
			}
			tx.Rollback()
			log.Fatal(err)
		}

		// Commit this question's transaction
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}

		// Log successful insert
		fmt.Printf("Imported question #%d with DB ID %d\n", i+1, questionID)
	}

	fmt.Println("Successfully imported questions!")
}
