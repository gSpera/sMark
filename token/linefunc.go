package token

import (
	"fmt"
	"strings"
)

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

func (l HeaderLine) String() string {
	return strings.Repeat(string(TypeHeader), 10)
}
func (l DivisorLine) String() string {
	return strings.Repeat(string(TypeLess), 10)
}
func (l LessLine) String() string {
	return strings.Repeat(string(TypeLess), l.Length)
}
func (l EqualLine) String() string {
	return strings.Repeat(string(TypeEqual), l.Length)
}
func (l ListLine) String() string {
	return strings.Repeat("\t", l.Indentation) + l.Text.String()
}
func (l CodeLine) String() string {
	return fmt.Sprintf("[%s]", l.Lang)
}
