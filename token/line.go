package token

import (
	"fmt"
)

//LineToken is a special interface that indicates a token
type LineToken interface{ IsToken() }

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
}

//String returns a string of the line with all tokens
func (t LineContainer) String() string {
	var str string
	for i := 0; i < t.Indentation; i++ {
		str += "\t"
	}

	str += t.StringNoTab()

	return str
}

//StringNoTab returns a string of the line with all tokens without tab
func (t LineContainer) StringNoTab() string {
	var str string

	for _, t := range t.Tokens {
		switch tt := t.(type) {
		case SimpleToken:
			str += string(tt.Char())
		case TextToken:
			str += tt.String()
		case CheckBoxToken:
			str += fmt.Sprintf("[%c]", tt.Char)
		default:
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", t, t))

		}
	}

	return str
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
	Indentation uint
}
