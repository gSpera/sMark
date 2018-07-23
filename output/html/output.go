package output

import (
	"bytes"
	"eNote/token"
	eNote "eNote/utils"
	"fmt"
	"html"
	"html/template"
	"log"
	"os"
	"strings"

	//Lexers for highlight
	_ "github.com/johnsto/go-highlight/lexers"
)

const maxMarkup = 255

//Titles:
//H1: Indent Title
//H3: Indent Subtitle
//H2: Title
//H4: Subtitle

//ToString is a simple output enging with a simple HTML writer
func ToString(paragraphs []token.ParagraphToken, options eNote.Options) []byte {
	var outTemplate *template.Template
	var err error
	if options.Bool["OnlyBody"] {
		outTemplate, err = template.New("Only Body").Parse(`{{.Body}}`)
	} else {
		outTemplate = template.New("HTML Output")
		tmpl, err := Asset("template.tmpl")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Output Engine: Could not get template: %v\n", err)
			os.Exit(1)
		}
		outTemplate.Parse(string(tmpl))
	}

	if err != nil {
		fmt.Println(err)
		panic("Output Engine: template is not valid")
	}

	body := ""
	// bold := false
	// italic := false
	alignMap := map[int]string{
		0: "align-left",
		1: "align-center",
		2: "align-right",
	}

	for _, p := range paragraphs {
		switch pp := p.(type) {
		case token.TitleParagraph:
			var tag string
			switch pp.Indentation {
			case 0:
				tag = "<h2>%s</h2>"
			default:
				tag = "<h1>%s</h1>"
			}
			body += fmt.Sprintf(tag, html.EscapeString(pp.Text))
			continue
		case token.DivisorParagraph:
			body += "<hr>"
			continue
		case token.TextParagraph:
			var paragraph string
			for _, line := range pp.Lines {
				paragraph += fromLineContainer(line)
				if options.Bool["NewLine"] {
					paragraph += "<br>\n"
				}
			}
			if paragraph == "" || paragraph == "<br>\n" {
				continue
			}
			body += fmt.Sprintf("<p class=\"%s\">%s</p>", alignMap[pp.Indentation], paragraph)
		case token.SubtitleParagraph:
			switch pp.Indentation {
			case 0:
				body += fmt.Sprintf("<h4>%s</h4>", html.EscapeString(pp.Text))
			default:
				body += fmt.Sprintf("<h3>%s</h3>", html.EscapeString(pp.Text))
			}
		case token.ListParagraph:
			currentIndentation := 0
			for _, item := range pp.Items {
				if item.Indentation > currentIndentation {
					for i := currentIndentation; i < item.Indentation; i++ {
						body += "<ul>"
					}
					currentIndentation = item.Indentation
				} else if item.Indentation < currentIndentation {
					for i := currentIndentation; i > item.Indentation; i-- {
						body += "</ul>"
					}
					currentIndentation = item.Indentation
				}
				body += fmt.Sprintf("<li>%s</li>", fromLineContainer(item.Text))
			}
			for i := currentIndentation; i > 0; i-- {
				body += "</ul>"
			}
		case token.CodeParagraph:
			body += fmt.Sprintf("<pre><code class=\"%s\">", pp.Lang)
			for _, line := range pp.Text.Lines {
				body += strings.Repeat("\t", line.Indentation) + fromLineContainer(line) + "\n"
			}
			body += "</code></pre>"
		default:
			log.Printf("Found Unknown Paragraph: %T{%+v}\n", p, p)
			panic("Paragraph Type not supported")
		}
	}

	var out bytes.Buffer
	err = outTemplate.Execute(&out, struct {
		Body    template.HTML
		Options eNote.Options
	}{template.HTML(body), options})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not execute template:", err)
	}
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

func fromLineContainer(line token.LineContainer) string {
	var str string
	for _, text := range line.Tokens {
		switch tt := text.(type) {
		case token.TextToken:
			if tt.Bold {
				str += "<b>"
			}
			if tt.Italic {
				str += "<i>"
			}
			if tt.Strike {
				str += "<s>"
			}
			str += tt.Text
			if tt.Bold {
				str += "</b>"
			}
			if tt.Italic {
				str += "</i>"
			}
			if tt.Strike {
				str += "</s>"
			}
		case token.CheckBoxToken:
			var data string
			if tt.Char != ' ' {
				data += "checked"
			}
			str += fmt.Sprintf("<input type=\"checkbox\" %s/>", data)
		default:
			if st, ok := tt.(token.SimpleToken); ok {
				str += string(st.Char())
				continue
			}
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", tt, tt))
		}
	}

	return str
}
