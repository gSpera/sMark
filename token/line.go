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
//TODO: Should be Tokens []TextToken??
type LineContainer struct {
	LineToken
	LineState
	Tokens []Token
}

//Type return the type of the Token
func (t LineContainer) Type() Type { return TypeTokenLine }

//String returns a string of the line with all tokens
//This method is deprecated
func (t LineContainer) String() string {
	var str string
	for i := 0; i < t.LineState.Indentation; i++ {
		str += "\t"
	}
	for _, t := range t.Tokens {
		if tt, ok := t.(SimpleToken); ok {
			str += string(tt.Char())
		} else {
			str += fmt.Sprintf("%v", t)
		}
	}
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
		default:
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", t, t))

		}
	}

	return str
}

//DebugString is used for ast outputting
func (t LineContainer) DebugString() string {
	return fmt.Sprintf("<LineToken: {%+v}%s>\n", t.LineState, t.Tokens)
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
