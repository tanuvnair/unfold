package parser_test

import (
	"strings"
	"testing"

	"github.com/tanuvnair/unfold/internal/parser"
)

func TestResolveColumnIndex_HitCaseInsensitive(t *testing.T) {
	header := []string{"Sl. No.", " Transaction Date ", "Description", "Amount"}
	idx, err := parser.ResolveColumnIndex(header, "description")
	if err != nil {
		t.Fatalf("ResolveColumnIndex: %v", err)
	}
	if idx != 2 {
		t.Fatalf("idx = %d, want 2", idx)
	}
}

func TestResolveColumnIndex_MissListsAvailable(t *testing.T) {
	header := []string{"Date", "Narration", "Amount"}
	_, err := parser.ResolveColumnIndex(header, "Description")
	if err == nil {
		t.Fatal("expected error for missing column")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Description") || !strings.Contains(msg, "Narration") {
		t.Fatalf("error should list wanted + available columns, got: %v", err)
	}
}

func TestResolveColumnIndex_EmptyName(t *testing.T) {
	_, err := parser.ResolveColumnIndex([]string{"Description"}, "  ")
	if err == nil {
		t.Fatal("expected error for empty column name")
	}
}
