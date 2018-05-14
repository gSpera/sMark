package output

import (
	"bytes"
	"eNote/token"
	eNote "eNote/utils"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
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
	bold := false
	italic := false
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
			body += fmt.Sprintf(tag, pp.Text)
			continue
		case token.DivisorParagraph:
			body += "<hr>"
			continue
		case token.TextParagraph:
			body += fmt.Sprintf("<p class=\"%s\">", alignMap[pp.Indentation])
			for _, line := range pp.Lines {
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
			body += "</p>\n"
		case token.SubtitleParagraph:
			log.Println("Found Subtitle:", pp.Text)
			switch pp.Indentation {
			case 0:
				body += fmt.Sprintf("<h4>%s</h4>", pp.Text)
			default:
				body += fmt.Sprintf("<h3>%s</h3>", pp.Text)
			}
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
