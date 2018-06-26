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
	if *options.OnlyBody {
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
			fmt.Println("TitleParagraph Indentation:", pp.Indentation)
			body += fmt.Sprintf(tag, html.EscapeString(pp.Text))
			continue
		case token.DivisorParagraph:
			body += "<hr>"
			continue
		case token.TextParagraph:
			var paragraph string
			for _, line := range pp.Lines {
				for _, text := range line.Tokens {
					switch tt := text.(type) {
					case token.TextToken:
						if tt.Bold {
							paragraph += "<b>"
						}
						if tt.Italic {
							paragraph += "<i>"
						}
						if tt.Strike {
							paragraph += "<s>"
						}
						paragraph += tt.Text
						if tt.Bold {
							paragraph += "</b>"
						}
						if tt.Italic {
							paragraph += "</i>"
						}
						if tt.Strike {
							paragraph += "</s>"
						}
					default:
						if st, ok := tt.(token.SimpleToken); ok {
							paragraph += string(st.Char())
							continue
						}
						panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", tt, tt))
					}
				}
				if *options.NewLine {
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
			body += "<ul>"

			for _, item := range pp.Items {
				body += fmt.Sprintf("<li>%s</li>", html.EscapeString(item.Text.String()))
			}

			body += "</ul>"
		default:
			log.Printf("Found Unknown Paragraph: %T{%+v}\n", p, p)
			panic("Paragraph Type not supported")
		}
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
