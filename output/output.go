package output

import (
	"eNote/token"
	"fmt"
)

//ToString oputput the list of tokens to string
func ToString(tokenList []token.Token) string {
	str := ""
	bold := false
	italic := false

	for _, _token := range tokenList {
		switch _token.(type) {
		case token.ParagraphToken:
		default:
			fmt.Printf("Type: %T\n", _token)
			panic("Top Level Token not supported")
		}

		paragraph := _token.(token.ParagraphToken)
		// str += fmt.Sprintf("<PARAGRAPH %d>\n", len(paragraph.Lines))
		fmt.Println("Output: Paragraph")
		for _, line := range paragraph.Lines {
			fmt.Println("Output: Line")
			// fmt.Printf("%T{%v}\n", line, line)
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
		}
		str += "<NEWLINE>\n"
	}

	return str
}
