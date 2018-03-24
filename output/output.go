package output

import (
	"eNote/token"
)

//ToString oputput the list of tokens to string
func ToString(tokenList []token.Token) string {
	str := ""
	bold := false
	italic := false

	for _, _line := range tokenList {
		switch _line.(type) {
		case token.LineToken:
		default:
			panic("Top Level Token not supported")
		}
		line := _line.(token.LineToken)

		for i := uint(0); i < line.Indentation; i++ {
			str += "<TAB>"
		}

		for _, t := range line.Tokens {
			switch t.(type) {
			case token.BoldToken:
				if !bold {
					str += "<BOLD"
				} else {
					str += ">"
				}
				bold = !bold
			case token.ItalicToken:
				if !italic {
					str += "<ITALIC"
				} else {
					str += ">"
				}
				italic = !italic
			default:
				str += t.String()
			}
		}
		str += "\n"
	}

	return str
}
