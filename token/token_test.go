package token

import (
	"strconv"
	"testing"
)

type token struct {
	Pos Position
	Typ Type
	Lit string
}

func (t *token) Position() Position {
	return t.Pos
}

func (t *token) Type() Type {
	return t.Typ
}
func (t *token) Literal() string {
	return t.Lit
}

func TestString(t *testing.T) {

	tks := []Token{
		&token{Lit: "a", Pos: Position{Filename: "", Line: 1, Column: 1}},
		&token{Lit: "中", Pos: Position{Filename: "", Line: 1, Column: 2}},
		&token{Lit: "♥", Pos: Position{Filename: "", Line: 1, Column: 3}},
		&token{Lit: "\n", Typ: NEWLINE, Pos: Position{Filename: "", Line: 1, Column: 4}},
		&token{Lit: "\n", Typ: NEWLINE, Pos: Position{Filename: "", Line: 2, Column: 1}},
		&token{Lit: "a", Pos: Position{Filename: "", Line: 3, Column: 1}},
		&token{Lit: "b", Pos: Position{Filename: "", Line: 4, Column: 1}},
		&token{Lit: "\n", Typ: NEWLINE, Pos: Position{Filename: "", Line: 4, Column: 2}},
		&token{Lit: "a", Pos: Position{Filename: "", Line: 5, Column: 1}},
		&token{Lit: "b", Pos: Position{Filename: "", Line: 6, Column: 1}},
	}

	got := String(tks)
	want := "a中♥\n\na\nb\na\nb"
	if want != got {
		t.Errorf("want %s got %s\n", strconv.QuoteToGraphic(want), strconv.QuoteToGraphic(got))
	}
}

func TestRelativeString(t *testing.T) {
	startLine := 10

	tks := []Token{
		&token{Lit: "a", Pos: Position{Filename: "", Line: 1 + startLine, Column: 1}},
		&token{Lit: "中", Pos: Position{Filename: "", Line: 1 + startLine, Column: 2}},
		&token{Lit: "♥", Pos: Position{Filename: "", Line: 1 + startLine, Column: 3}},
		&token{Lit: "\n", Typ: NEWLINE, Pos: Position{Filename: "", Line: 1 + startLine, Column: 4}},
		&token{Lit: "\n", Typ: NEWLINE, Pos: Position{Filename: "", Line: 2 + startLine, Column: 1}},
		&token{Lit: "a", Pos: Position{Filename: "", Line: 3 + startLine, Column: 1}},
		&token{Lit: "b", Pos: Position{Filename: "", Line: 4 + startLine, Column: 1}},
		&token{Lit: "\n", Typ: NEWLINE, Pos: Position{Filename: "", Line: 4 + startLine, Column: 2}},
		&token{Lit: "a", Pos: Position{Filename: "", Line: 5 + startLine, Column: 1}},
		&token{Lit: "b", Pos: Position{Filename: "", Line: 6 + startLine, Column: 1}},
	}

	got := RelativeString(tks)
	want := "a中♥\n\na\nb\na\nb"
	if want != got {
		t.Errorf("want %s got %s\n", strconv.QuoteToGraphic(want), strconv.QuoteToGraphic(got))
	}
}
