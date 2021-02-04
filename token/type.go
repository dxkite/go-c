package token

import "strconv"

type Type int

const (
	ILLEGAL Type = iota
	EOF
	INT   // int number
	FLOAT //  float number
	IDENT
	CHAR
	STRING
	NEWLINE
	WHITESPACE
	PUNCTUATOR
	KEYWORD
	TEXT
)

var tokenName = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	INT:        "INT",
	FLOAT:      "FLOAT",
	IDENT:      "IDENT",
	CHAR:       "CHAR",
	STRING:     "STRING",
	NEWLINE:    "NEWLINE",
	WHITESPACE: "WHITESPACE",
	PUNCTUATOR: "PUNCTUATOR",
	KEYWORD:    "KEYWORD",
	TEXT:       "TEXT",
}

var nameToken = map[string]Type{
	"ILLEGAL":    ILLEGAL,
	"EOF":        EOF,
	"INT":        INT,
	"FLOAT":      FLOAT,
	"IDENT":      IDENT,
	"CHAR":       CHAR,
	"STRING":     STRING,
	"NEWLINE":    NEWLINE,
	"WHITESPACE": WHITESPACE,
	"PUNCTUATOR": PUNCTUATOR,
	"KEYWORD":    KEYWORD,
	"TEXT":       TEXT,
}

type Name string

func (tok Name) Token() Type {
	if v, ok := nameToken[string(tok)]; ok {
		return v
	}
	return 0
}

func (tok Type) String() string {
	s := ""
	if 0 <= tok && int(tok) < len(tokenName) {
		s = tokenName[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

func (f Type) MarshalJSON() ([]byte, error) {
	typ := f.String()
	return []byte(strconv.QuoteToGraphic(typ)), nil
}

func (f *Type) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	*f = Name(str).Token()
	return nil
}
