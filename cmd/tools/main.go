package main

import (
	"max-odyssey-backend/utils"
	"os"
)

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
