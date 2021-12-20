package token

// token
type Token interface {
	Position() Position
	Type() Type
	Literal() string
}
