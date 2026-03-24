package parser_test

import (
	"os"
	"testing"

	"github.com/ElyasAsmad/everestengineering2/internal/model"
	"github.com/ElyasAsmad/everestengineering2/internal/parser"
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

func TestParseOffersCSV_NormalCSV(t *testing.T) {
	content := `code,discount,distance,weight
OFR001,20,d < 100,w < 70
OFR002,10,50 <= d <= 150,70 <= w <= 200
OFR003,5,d > 150,w > 200`

	tmpFile := createTempCSV(t, content)

	offers, err := parser.ParseOffersCSV(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(offers) != 3 {
		t.Fatalf("expected 3 offers, got %d", len(offers))
	}

	expected := []*model.Offer{
		{
			Code:       "OFR001",
			Discount:   20,
			Constraint: "d < 100 && w < 70",
		},
		{
			Code:       "OFR002",
			Discount:   10,
			Constraint: "50 <= d <= 150 && 70 <= w <= 200",
		},
		{
			Code:       "OFR003",
			Discount:   5,
			Constraint: "d > 150 && w > 200",
		},
	}

	for i := range expected {
		if offers[i].Code != expected[i].Code {
			t.Errorf("expected code %s, got %s", expected[i].Code, offers[i].Code)
		}
		if offers[i].Discount != expected[i].Discount {
			t.Errorf("expected discount %.2f, got %.2f", expected[i].Discount, offers[i].Discount)
		}
		if offers[i].Constraint != expected[i].Constraint {
			t.Errorf("expected constraint '%s', got '%s'", expected[i].Constraint, offers[i].Constraint)
		}
	}
}

func TestParseOffersCSV_CSVWithLeadingSpace(t *testing.T) {
	content := `code,discount,distance,weight
OFR001,20, d < 100, w < 70
OFR002,10, 50 <= d <= 150, 70 <= w <= 200
OFR003,5, d > 150, w > 200`

	tmpFile := createTempCSV(t, content)

	offers, err := parser.ParseOffersCSV(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(offers) != 3 {
		t.Fatalf("expected 3 offers, got %d", len(offers))
	}

	expected := []*model.Offer{
		{
			Code:       "OFR001",
			Discount:   20,
			Constraint: "d < 100 && w < 70",
		},
		{
			Code:       "OFR002",
			Discount:   10,
			Constraint: "50 <= d <= 150 && 70 <= w <= 200",
		},
		{
			Code:       "OFR003",
			Discount:   5,
			Constraint: "d > 150 && w > 200",
		},
	}

	for i := range expected {
		if offers[i].Code != expected[i].Code {
			t.Errorf("expected code %s, got %s", expected[i].Code, offers[i].Code)
		}
		if offers[i].Discount != expected[i].Discount {
			t.Errorf("expected discount %.2f, got %.2f", expected[i].Discount, offers[i].Discount)
		}
		if offers[i].Constraint != expected[i].Constraint {
			t.Errorf("expected constraint '%s', got '%s'", expected[i].Constraint, offers[i].Constraint)
		}
	}
}

func TestParseOffersCSV_InvalidCSVNotEnoughColumns(t *testing.T) {
	content := `code,discount,distance,weight
OFR001,20,d < 100`

	tmpFile := createTempCSV(t, content)

	_, err := parser.ParseOffersCSV(tmpFile)
	if err == nil {
		t.Fatal("expected error for invalid CSV")
	}
}

func TestParseOffersCSV_NoRows(t *testing.T) {
	content := `code,discount,distance,weight`

	tmpFile := createTempCSV(t, content)

	offers, err := parser.ParseOffersCSV(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(offers) != 0 {
		t.Fatalf("expected 0 offers, got %d", len(offers))
	}
}

func TestParseOffersCSV_InvalidCSVNonNumericDiscount(t *testing.T) {
	content := `code,discount,distance,weight
OFR001,twenty,d < 100,w < 70`

	tmpFile := createTempCSV(t, content)

	_, err := parser.ParseOffersCSV(tmpFile)
	if err == nil {
		t.Fatal("expected error for invalid CSV")
	}
}

func TestParseOffersCSV_NoConstraints(t *testing.T) {
	content := `code,discount,distance,weight
OFR001,20,,`

	tmpFile := createTempCSV(t, content)

	offers, err := parser.ParseOffersCSV(tmpFile)
	if err != nil {
		t.Fatal("expected error for invalid CSV")
	}

	if len(offers) != 1 {
		t.Fatalf("expected 1 offer, got %d", len(offers))
	}

	expected := []*model.Offer{
		{
			Code:       "OFR001",
			Discount:   20,
			Constraint: "",
		},
	}

	for i := range expected {
		if offers[i].Constraint != expected[i].Constraint {
			t.Errorf("expected constraint '%s', got '%s'", expected[i].Constraint, offers[i].Constraint)
		}
	}
}

func TestParseOffersCSV_SingleConstraint(t *testing.T) {
	content := `code,discount,distance,weight
OFR001,20,,w < 70
OFR002,10,d < 100,`

	tmpFile := createTempCSV(t, content)

	offers, err := parser.ParseOffersCSV(tmpFile)
	if err != nil {
		t.Fatal("expected error for invalid CSV")
	}

	if len(offers) != 2 {
		t.Fatalf("expected 2 offers, got %d", len(offers))
	}

	expected := []*model.Offer{
		{
			Code:       "OFR001",
			Discount:   20,
			Constraint: "w < 70",
		},
		{
			Code:       "OFR002",
			Discount:   10,
			Constraint: "d < 100",
		},
	}

	for i := range expected {
		if offers[i].Constraint != expected[i].Constraint {
			t.Errorf("expected constraint '%s', got '%s'", expected[i].Constraint, offers[i].Constraint)
		}
	}
}
