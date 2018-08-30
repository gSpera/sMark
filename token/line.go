package token

import (
	"fmt"
)

//LineToken is a special interface that indicates a line
type LineToken interface {
	IsToken()
	fmt.Stringer
}

//HeaderLine rapresent a header line
type HeaderLine struct {
	LineToken
	Tokens []Token
}

//LineContainer is a token which rappresent a list of Tokens with some attributes
type LineContainer struct {
	LineToken
	Indentation int
	Tokens      []Token

	//Attribute
	Quote bool
}

//DivisorLine rapresent a header line
type DivisorLine struct {
	LineToken
}

//LessLine rapresent a line containing only LessToken
type LessLine struct {
	LineToken
	Indentation int
	Length      int
}

//EqualLine rapresent a line containing only EqualToken
//EqualLine is used for Titles(Main and not)
type EqualLine struct {
	LineToken
	Indentation int
	Length      int
}

//ListLine rapresent an element of a list
type ListLine struct {
	LineToken
	Text        LineContainer
	Indentation int
}

//CodeLine is a line containing the language name for CodeParagraph
type CodeLine struct {
	LineToken
	Lang string
}
