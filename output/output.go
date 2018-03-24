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

//ToString is a simple output enging with a simple HTML writer
func ToString(paragraphs []token.ParagraphToken) []byte {
	str := `<html>
	<head>
	<style>
		p.align-left {text-align:left;}
		p.align-center {text-align: center;}
		p.align-right {text-align: right;}
	</style>
	</head>
	<body>`
	bold := false
	italic := false
	alignMap := map[int]string{
		0: "align-left",
		1: "align-center",
		2: "align-right",
	}

	for _, p := range paragraphs {
		str += fmt.Sprintf("<p class=\"%s\">", alignMap[p.Indentation])
		for _, line := range p.Lines {
			for _, tok := range line.Tokens {

				switch tok.(type) {
				case token.BoldToken:
					fmt.Println("Bold")
					bold = !bold
					if bold {
						str += "<b>"
					} else {
						str += "</b> "
					}
					continue
				case token.ItalicToken:
					fmt.Println("Italic")
					italic = !italic
					if italic {
						str += "<i>"
					} else {
						str += "</i> "
					}
					continue
				}

				fmt.Printf("Adding Text: %s, Bold: %v, Italic: %v\n", tok.String(), bold, italic)
				str += tok.String()
				if bold {
					fmt.Println("Apply bold")
					fmt.Println(tok.String())
				}
				if italic {
					fmt.Println("Apply Italic")
					fmt.Println(tok.String())
				}

			}
			str += "<br>\n"
		}
		str += "</p>\n"
	}

	str += "</body>\n</html>"
	return []byte(str)
}
