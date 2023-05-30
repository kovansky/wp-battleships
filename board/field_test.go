package board_test

import (
	"github.com/kovansky/wp-battleships/board"
	"testing"
)

func TestField_Numeric(t *testing.T) {
	type tableData struct {
		name     string
		input    string
		expected uint8
		wantErr  bool
	}

	table := []tableData{
		{"Single letter A", "A", 0, true},
		{"Two letter A", "AA", 0, true},
		{"Single letter G", "G", 0, true},
		{"Correct A1", "A1", 0, false},
		{"Correct B1", "B1", 10, false},
		{"Correct C3", "C3", 22, false},
		{"Correct J10", "J10", 99, false},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			f, err := board.NewField(tt.input)
			if err != nil && !tt.wantErr {
				t.Fatalf("Received unexpected error: %v", err)
			} else if err != nil && tt.wantErr {
				return
			}

			got := f.Numeric()

			if got != tt.expected {
				t.Fatalf("Incorrect numeric value; expected: %d, got: %d", tt.expected, got)
			}
		})
	}
}

func TestField_IsCorner(t *testing.T) {
	type tableData struct {
		name     string
		field    string
		target   string
		expected bool
	}

	table := []tableData{
		{"A1 - A2 (not corner)", "A1", "A2", false},
		{"D4 - C3 (corner)", "D4", "C3", true},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			f, err := board.NewField(tt.field)
			if err != nil {
				t.Fatalf("Received unexpected error: %v", err)
			}

			got := f.IsCorner(tt.target)
			
			if got != tt.expected {
				t.Fatalf("Incorrect numeric value; expected: %t, got: %t", tt.expected, got)
			}
		})
	}
}
