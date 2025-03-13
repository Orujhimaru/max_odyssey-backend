package main

import (
	"max-odyssey-backend/utils"
	"os"
)

// in order to delete the db, run the command:
// # First delete from answer_choices
// docker exec -it max_odyssey-backend-postgres-1 psql -U satapp -d sat_tracker -c "DELETE FROM answer_choices;"

// Then delete from questions
// docker exec -it max_odyssey-backend-postgres-1 psql -U satapp -d sat_tracker -c "DELETE FROM questions;"

// go run cmd/tools/main.go import  using the command, imports the questions,
// from utils/import_questions.go
func main() {
	switch os.Args[1] {
	case "import":
		utils.ImportQuestions()
	case "add":
		utils.AddQuestion()
	default:
		println("Available commands: import, add")
	}
}
