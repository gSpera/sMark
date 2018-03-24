package output

import (
	"eNote/token"
	"fmt"
)

//DebugToString oputput the list of tokens to string
func DebugToString(tokenList []token.ParagraphToken) string {
	str := ""
	bold := false
	italic := false

	for _, paragraph := range tokenList {
		// str += fmt.Sprintf("<PARAGRAPH %d>\n", len(paragraph.Lines))
		fmt.Println("Output: Paragraph")
		str += fmt.Sprintf("<PARAGRAPH %d>\n", paragraph.Indentation)
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
						str += "<BOLD "
					} else {
						str += ">"
					}
					bold = !bold
				case token.ItalicToken:
					if !italic {
						str += "<ITALIC "
					} else {
						str += ">"
					}
					italic = !italic
				default:
					str += t.String()
				}
			}
			str += "<NEWLINE>\n"
		}
	}

	return str
}
