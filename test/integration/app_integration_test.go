//go:build integration

package integration_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ElyasAsmad/everestengineering2/internal/app"
)

func createTempCSV(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "offers_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	return tmpFile.Name()
}

func TestFullFlow(t *testing.T) {
	offers := `code,discount,distance,weight
OFR001,10,d < 200, 70 <= w <= 200
OFR002,7,50 <= d <= 150, 100 <= w <= 250
OFR003,5,50 <= d <= 250, 10 <= w <= 150`

	input := `100 5
PKG1 50 30 OFR001
PKG2 75 125 OFR008
PKG3 175 100 OFR003
PKG4 110 60 OFR002
PKG5 155 95 NA
2 70 200`

	expected := `PKG1 0 750 3.98
PKG2 0 1475 1.78
PKG3 0 2350 1.42
PKG4 105 1395 0.85
PKG5 0 2125 4.19`

	tmpFile := createTempCSV(t, offers)

	reader := strings.NewReader(input)

	result, err := app.Run(reader, tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestRun_MissingInput(t *testing.T) {
	offers := `code,discount,distance,weight
OFR001,10,d < 200, 70 <= w <= 200
OFR002,7,50 <= d <= 150, 100 <= w <= 250
OFR003,5,50 <= d <= 250, 10 <= w <= 150`

	tmpFile := createTempCSV(t, offers)

	reader := strings.NewReader("")

	_, err := app.Run(reader, tmpFile)
	if err == nil {
		t.Fatal("Expected error for invalid input format, got nil")
	}
}

func TestRun_InvalidCSVFile(t *testing.T) {
	input := `100 5
PKG1 50 30 OFR001
PKG2 75 125 OFR008
PKG3 175 100 OFR003
PKG4 110 60 OFR002
PKG5 155 95 NA
2 70 200`

	_, err := app.Run(strings.NewReader(input), "nonexistent.csv")
	if err == nil {
		t.Fatal("Expected error for nonexistent CSV file, got nil")
	}
}

func TestRun_InvalidInputFormat(t *testing.T) {
	offers := `code,discount,distance,weight
OFR001,10,d < 200, 70 <= w <= 200
OFR002,7,50 <= d <= 150, 100 <= w <= 250
OFR003,5,50 <= d <= 250, 10 <= w <= 150`

	fileName := createTempCSV(t, offers)

	_, err := app.Run(strings.NewReader("invalid input"), fileName)
	if err == nil {
		t.Fatal("Expected error for invalid input format, got nil")
	}
}

func TestRun_MissingVehicleDetails(t *testing.T) {
	offers := `code,discount,distance,weight
OFR001,10,d < 200, 70 <= w <= 200
OFR002,7,50 <= d <= 150, 100 <= w <= 250
OFR003,5,50 <= d <= 250, 10 <= w <= 150`

	// last line should contain vehicle details, but it's missing
	input := `100 5
PKG1 50 30 OFR001
PKG2 75 125 OFR008
PKG3 175 100 OFR003
PKG4 110 60 OFR002
PKG5 155 95 NA`

	tmpFile := createTempCSV(t, offers)

	reader := strings.NewReader(input)

	_, err := app.Run(reader, tmpFile)
	if err == nil {
		t.Fatal("Expected error for missing vehicle details, got nil")
	}
}
