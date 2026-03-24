package main

import (
	"fmt"
	"os"

	"github.com/ElyasAsmad/everestengineering2/internal/app"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	result := app.Run(inputFile)

	fmt.Print(result)
}
