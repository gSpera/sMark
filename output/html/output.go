package output

import (
	"bytes"
	"eNote/token"
	eNote "eNote/utils"
	"fmt"
	"html/template"
	"path"
	"strings"
)

const maxMarkup = 255

//DebugToString oputput the list of tokens to string
//You should not use this function.
func DebugToString(tokenList []token.ParagraphToken, options eNote.Options) []byte {
	str := ""
	bold := false
	italic := false

	for _, paragraph := range tokenList {
		paragraph := paragraph.(token.TextParagraph)
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

	return []byte(str)
}

//ToString is a simple output enging with a simple HTML writer
func ToString(paragraphs []token.ParagraphToken, options eNote.Options) []byte {
	var outTemplate *template.Template
	var err error
	if *options.OnlyBody {
		outTemplate, err = template.New("Only Body").Parse(`{{.Body}}`)
	} else {
		outTemplate, err = template.ParseFiles(path.Join("output", "html", "template.tmpl"))
	}
	if err != nil {
		fmt.Println(err)
		panic("Output Engine: template is not valid")
	}

	body := ""
	bold := false
	italic := false
	alignMap := map[int]string{
		0: "align-left",
		1: "align-center",
		2: "align-right",
	}

	for _, p := range paragraphs {
		switch p.(type) {
		case token.TextParagraph:
			p := p.(token.TextParagraph)
			body += fmt.Sprintf("<p class=\"%s\">", alignMap[p.Indentation])
			for _, line := range p.Lines {
				for i, tok := range line.Tokens {
					switch tok.(type) {
					case token.BoldToken:
						if distance := findToken(line, i, token.TypeBold); !bold && distance > maxMarkup || distance == -1 {
							break
						}

						fmt.Println("Bold")
						bold = !bold
						if bold {
							body += "<b>"
						} else {
							body += "</b> "
						}
						continue
					case token.ItalicToken:
						fmt.Println("Italic")
						italic = !italic
						if italic {
							body += "<i>"
						} else {
							body += "</i> "
						}
						continue
					case token.TabToken:
						fmt.Println("TAB")
						body += strings.Repeat("&nbsp;", int(*options.TabWidth))
						continue
					}

					fmt.Printf("Adding Text: %s, Bold: %v, Italic: %v\n", tok.String(), bold, italic)
					body += tok.String()
					if bold {
						fmt.Println("Apply bold")
						fmt.Println(tok.String())
					}
					if italic {
						fmt.Println("Apply Italic")
						fmt.Println(tok.String())
					}

				}

				if *options.NewLine {
					body += "<br>\n"
				}
			}
		}
		body += "</p>\n"
	}

	var out bytes.Buffer
	outTemplate.Execute(&out, struct {
		Body    template.HTML
		Options eNote.OptionsTemplate
	}{template.HTML(body), options.ToTemplate()})
	return out.Bytes()
}

func findToken(line token.LineContainer, start int, t token.Type) int {
	for i := start; i < len(line.Tokens); i++ {
		switch line.Tokens[i].Type() {
		case t:
			return i
		}
	}
	return -1
}
