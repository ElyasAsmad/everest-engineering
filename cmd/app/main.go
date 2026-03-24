package main

import (
	"fmt"
	"os"

	"github.com/ElyasAsmad/everestengineering2/internal/app"
	"github.com/ElyasAsmad/everestengineering2/internal/logger"
)

func main() {
	logger := logger.NewLogger()

	if len(os.Args) != 2 {
		logger.Fatal("Usage: go run main.go <input_file.csv>")
	}

	inputFile := os.Args[1]
	result, err := app.Run(os.Stdin, inputFile)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Print(result)
}
