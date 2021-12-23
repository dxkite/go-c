package token

import (
	"fmt"
	"strconv"
)

// token
type Token interface {
	Position() Position
	Type() Type
	Literal() string
}

func String(t Token) string {
	position := t.Position()
	return fmt.Sprintf("%s<%s@%s>", strconv.QuoteToGraphic(t.Literal()), t.Type(), position.String())
}
