package utils_test

import (
	"testing"

	"github.com/ElyasAsmad/everestengineering2/internal/utils"
)

type TestCase struct {
	input float64
	want  int
}

func TestRound(t *testing.T) {
	tests := []TestCase{
		{input: 1.4, want: 1},
		{input: 1.5, want: 2},
		{input: -1.5, want: -2},
		{input: -1.4, want: -1},
	}

	for _, tc := range tests {
		got := utils.Round(tc.input)
		if got != tc.want {
			t.Errorf("Round(%v) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestToFixed(t *testing.T) {
	got := utils.ToFixed(3.14159, 2)
	want := 3.14

	if got != want {
		t.Errorf("ToFixed() = %v; want %v", got, want)
	}
}

func TestTruncate(t *testing.T) {
	got := utils.Truncate(1.239, 2)
	want := 1.23

	if got != want {
		t.Errorf("Truncate() = %v; want %v", got, want)
	}
}
