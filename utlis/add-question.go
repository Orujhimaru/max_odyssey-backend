package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5431 user=satapp password=satapp123 dbname=sat_tracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter subject ID: ")
	subjectIDStr, _ := reader.ReadString('\n')
	subjectID, err := strconv.Atoi(subjectIDStr[:len(subjectIDStr)-1])
	if err != nil {
		log.Fatalf("Invalid subject ID: %v", err)
	}

	fmt.Print("Enter question text: ")
	questionText, _ := reader.ReadString('\n')

	fmt.Print("Enter correct answer: ")
	correctAnswer, _ := reader.ReadString('\n')

	fmt.Print("Enter difficulty level (1-5): ")
	difficultyLevelStr, _ := reader.ReadString('\n')
	difficultyLevel, err := strconv.Atoi(difficultyLevelStr[:len(difficultyLevelStr)-1])
	if err != nil {
		log.Fatalf("Invalid difficulty level: %v", err)
	}

	fmt.Print("Enter explanation: ")
	explanation, _ := reader.ReadString('\n')

	_, err = db.Exec("INSERT INTO questions (subject_id, question_text, correct_answer, difficulty_level, explanation) VALUES ($1, $2, $3, $4, $5)",
		subjectID, questionText[:len(questionText)-1], correctAnswer[:len(correctAnswer)-1], difficultyLevel, explanation[:len(explanation)-1])
	if err != nil {
		log.Fatalf("Failed to insert question: %v", err)
	}

	fmt.Println("Question added successfully!")
}
