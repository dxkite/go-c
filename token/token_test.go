package token

import (
	"strconv"
	"testing"
)

func TestString(t *testing.T) {

	tks := []*Token{
		{Lit: "a", Position: Position{Filename: "", Line: 1, Column: 1}},
		{Lit: "中", Position: Position{Filename: "", Line: 1, Column: 2}},
		{Lit: "♥", Position: Position{Filename: "", Line: 1, Column: 3}},
		{Lit: "\n", Type: NEWLINE, Position: Position{Filename: "", Line: 1, Column: 4}},
		{Lit: "\n", Type: NEWLINE, Position: Position{Filename: "", Line: 2, Column: 1}},
		{Lit: "a", Position: Position{Filename: "", Line: 3, Column: 1}},
		{Lit: "b", Position: Position{Filename: "", Line: 4, Column: 1}},
		{Lit: "\n", Type: NEWLINE, Position: Position{Filename: "", Line: 4, Column: 2}},
		{Lit: "a", Position: Position{Filename: "", Line: 5, Column: 1}},
		{Lit: "b", Position: Position{Filename: "", Line: 6, Column: 1}},
	}

	got := String(tks)
	want := "a中♥\n\na\nb\na\nb"
	if want != got {
		t.Errorf("want %s got %s\n", strconv.QuoteToGraphic(want), strconv.QuoteToGraphic(got))
	}
}
