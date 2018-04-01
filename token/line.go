package token

import "fmt"

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
	LineState
	Tokens []Token
}

//Type return the type of the Token
func (t LineContainer) Type() Type { return TypeTokenLine }

func (t LineContainer) String() string {
	str := ""
	for i := uint(0); i < t.LineState.Indentation; i++ {
		str += "\t"
	}
	for _, t := range t.Tokens {
		str += fmt.Sprintf("%v", t)
	}
	return str + "\n"
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
}
