package htmlout

import (
	"bytes"
	"eNote/token"
	eNote "eNote/utils"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma"
	chtml "github.com/alecthomas/chroma/formatters/html"

	//Lexers for highlight
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const (
	//DefaultStyle is the default theme used when highlight code
	DefaultStyle = "vs"

	//StyleOption is the string options used to define the style
	StyleOption = "Style"

	//PrintLineOption is the bool options used to determine wheteever to print the line number
	//near code lines, default true
	PrintLineOption = "ShowCodeLines"

	//OnlyBodyOption is the Bool Option that specify to output only the body of the HTML Document
	OnlyBodyOption = "OnlyBody"

	//CustomTemplateOption is the String Option used to specify a custom template
	CustomTemplateOption = "TemplateFile"
)

//Titles:
//H1: Indent Title
//H3: Indent Subtitle
//H2: Title
//H4: Subtitle

//ToString is a simple output enging with a simple HTML writer
func ToString(paragraphs []token.ParagraphToken, options eNote.Options) []byte {
	outTemplate := template.New("eNote")
	var err error

	outTemplate = outTemplate.Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"safeCSS":  func(s string) template.CSS { return template.CSS(s) },
		"safeJS":   func(s string) template.JS { return template.JS(s) },
	})

	if options.Bool[OnlyBodyOption] {
		outTemplate, err = outTemplate.Parse(`{{.Body}}`)
	} else if tt, ok := options.String[CustomTemplateOption]; ok {
		log.Println("\t- Parsing:", tt)
		content, err := ioutil.ReadFile(tt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot read file:", tt)
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		outTemplate, err = outTemplate.Parse(string(content))
	} else {
		tmpl, err := Asset("template.tmpl")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Output Engine: Could not get template: %v\n", err)
			os.Exit(1)
		}
		outTemplate.Parse(string(tmpl))
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse template:", err)
		os.Exit(1)
	}

	head := HTMLNode{
		tag: "",
	}
	body := HTMLNode{
		tag: "",
	}

	alignMap := map[int]string{
		0: "align-left",
		1: "align-center",
		2: "align-right",
	}

	for _, p := range paragraphs {
		switch pp := p.(type) {
		case token.TitleParagraph:
			tag := "h1"
			if pp.Indentation == 0 {
				tag = "h2"
			}

			body.AddChildren(&HTMLNode{
				tag: tag,
				children: []Node{
					TextNode(pp.Text),
				},
			})
			continue
		case token.DivisorParagraph:
			body.AddChildren(&HTMLNode{
				tag:    "hr",
				single: true,
			})
			continue
		case token.TextParagraph:
			var quote bool
			paragraph := &HTMLNode{
				tag: "",
			}
			current := paragraph
			for _, line := range pp.Lines {
				if line.Quote && !quote {
					child := &HTMLNode{
						tag: "blockquote",
					}
					current.AddChildren(child)
					current = child
				} else if !line.Quote && quote {
					//Remove blockquote
					current = current.parent
				}
				quote = line.Quote
				current.AddChildrenNode(fromLineContainer(line))
				if options.Bool["NewLine"] {
					current.AddChildren(&HTMLNode{
						tag:    "br",
						single: true,
					})
				}
			}

			if len(paragraph.children) == 0 ||
				(len(paragraph.children) == 1 && paragraph.children[0].HTML(0) == "<br />") {
				continue
			}
			body.AddChildren(&HTMLNode{
				tag: "p",
				attrs: map[string]string{
					"class": alignMap[pp.Indentation],
				},
				children: []Node{
					paragraph,
				},
			})

		case token.SubtitleParagraph:
			tag := "h3"
			if pp.Indentation == 0 {
				tag = "h4"
			}

			body.AddChildren(&HTMLNode{
				tag: tag,
				children: []Node{
					TextNode(html.EscapeString(pp.Text)),
				},
			})
		case token.ListParagraph:
			currentIndentation := 1
			list := &HTMLNode{
				tag: "ul",
			}
			current := list

			for _, item := range pp.Items {
				if item.Indentation > currentIndentation {
					for i := currentIndentation; i < item.Indentation; i++ {
						child := &HTMLNode{
							tag: "ul",
						}
						current.AddChildren(child)
						current = child
					}
					currentIndentation = item.Indentation
				} else if item.Indentation < currentIndentation {
					for i := currentIndentation; i > item.Indentation; i-- {
						current = current.parent
					}
					currentIndentation = item.Indentation
				}
				current.AddChildren(&HTMLNode{
					tag: "li",
					children: []Node{
						fromLineContainer(item.Text),
					},
				})
			}
			body.AddChildren(list)
		case token.CodeParagraph:
			hd, bd := generateCode(pp, options)
			if hd != nil {
				head.AddChildren(hd)
			}
			body.AddChildrenNode(CodeNode(bd))
		default:
			log.Printf("Found Unknown Paragraph: %T{%+v}\n", p, p)
			panic("Paragraph Type not supported")
		}
	}

	var out bytes.Buffer
	options.String["_Head"] = head.HTML(0)

	err = outTemplate.Execute(&out, struct {
		Body    template.HTML
		Bool    map[string]bool
		String  map[string]string
		Generic map[string]interface{}
	}{template.HTML(body.HTML(0)), options.Bool, options.String, options.Generic})
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

func fromLineContainer(line token.LineContainer) *HTMLNode {
	root := &HTMLNode{
		tag: "",
	}

	for _, text := range line.Tokens {
		switch tt := text.(type) {
		case token.TextToken:
			textRoot := &HTMLNode{
				tag: "",
			}
			current := textRoot

			if tt.Bold {
				node := &HTMLNode{
					tag: "b",
				}
				current.AddChildren(node)
				current = node
			}
			if tt.Italic {
				node := &HTMLNode{
					tag: "i",
				}
				current.AddChildren(node)
				current = node
			}
			if tt.Strike {
				node := &HTMLNode{
					tag: "s",
				}
				current.AddChildren(node)
				current = node
			}
			if tt.Link != "" {
				node := &HTMLNode{
					tag: "a",
					attrs: map[string]string{
						"href": tt.Link,
					},
				}
				current.AddChildren(node)
				current = node
			}

			current.AddChildrenNode(TextNode(tt.Text))
			root.AddChildrenNode(textRoot.children[0])
		case token.CheckBoxToken:
			attrs := map[string]string{
				"type": "checkbox",
			}
			if tt.Char != ' ' {
				attrs["checked"] = "true"
			}
			root.AddChildren(&HTMLNode{
				tag:   "input",
				attrs: attrs,
			})
		case token.SimpleToken:
			root.AddChildrenNode(TextNode(string(tt.Char())))

		default:
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", tt, tt))
		}
	}

	return root
}

func generateCode(p token.CodeParagraph, options eNote.Options) (*HTMLNode, string) {
	lex := lexers.Get(p.Lang)
	if lex == nil {
		fmt.Fprintln(os.Stderr, "Cannot highlight lang:", p.Lang)
		lex = lexers.Fallback
	}
	lex = chroma.Coalesce(lex)

	stl := DefaultStyle
	if st, ok := options.String[StyleOption]; ok {
		stl = st
	}
	style := styles.Get(stl)
	if style.Name != stl {
		fmt.Fprintf(os.Stderr, "Style %s not found, using default\n", stl)
		style = styles.Get(DefaultStyle)
		if style.Name != DefaultStyle {
			fmt.Fprintf(os.Stderr, "Cannot get DefaultStyle %s, using FallBack\n", DefaultStyle)
			style = styles.Fallback
		}
	}

	fOpts := []chtml.Option{
		chtml.WithClasses(),
		chtml.ClassPrefix("hl-"),
	}
	if v, ok := options.Bool[PrintLineOption]; !ok || v {
		fOpts = append(fOpts,
			chtml.WithLineNumbers(),
			chtml.LineNumbersInTable(),
		)
	}
	formatter := chtml.New(fOpts...)

	if formatter == nil {
		fmt.Fprintln(os.Stderr, "Cannot create formatter")
	}

	highlight := bytes.NewBuffer([]byte{})
	var code string
	for _, l := range p.Text.Lines {
		code += l.String() + "\n"
	}

	iter, err := lex.Tokenise(nil, code)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not genereta highlight code")
		return nil, ""
	}

	err = formatter.Format(highlight, style, iter)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Cannot format highlight:", err)
		return nil, ""
	}

	css := bytes.NewBuffer([]byte{})
	formatter.WriteCSS(css, style)

	cssNode := &HTMLNode{
		tag: "style",
		children: []Node{
			TextNode(css.String()),
		},
	}
	out := highlight.String()
	return cssNode, out
}

func linecontainerToText(p token.TextParagraph) TextNode {
	var str string

	for _, line := range p.Lines {
		str += strings.Repeat("\t", line.Indentation) + fromLineContainer(line).HTML(0) + "\n"
	}
	return TextNode(str)
}
