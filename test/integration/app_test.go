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

	result := app.Run(reader, tmpFile)

	t.Log(result)

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
