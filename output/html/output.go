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
)

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
	} else if tt, ok := options.String["TemplateFile"]; ok {
		log.Println("\t- Parsing:", tt)
		content, err := ioutil.ReadFile(tt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot read file:", tt)
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		outTemplate, err = template.New("Custom Template").Parse(string(content))
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
		fmt.Fprintln(os.Stderr, "Cannot parse template:", err)
		os.Exit(1)
	}

	head := HtmlNode{
		tag: "",
	}
	body := HtmlNode{
		tag: "",
	}
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
			tag := "h1"
			if pp.Indentation == 0 {
				tag = "h2"
			}

			body.AddChildren(&HtmlNode{
				tag: tag,
				children: []Node{
					TextNode(pp.Text),
				},
			})
			continue
		case token.DivisorParagraph:
			body.AddChildren(&HtmlNode{
				tag:    "hr",
				single: true,
			})
			continue
		case token.TextParagraph:
			var quote bool
			paragraph := &HtmlNode{
				tag: "p",
			}
			current := paragraph
			for _, line := range pp.Lines {
				if line.Quote && !quote {
					child := &HtmlNode{
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
					current.AddChildren(&HtmlNode{
						tag:    "br",
						single: true,
					})
				}
			}

			if len(paragraph.children) == 0 ||
				(len(paragraph.children) == 1 && paragraph.children[0].HTML(0) == "<br />") {
				continue
			}
			body.AddChildren(&HtmlNode{
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

			body.AddChildren(&HtmlNode{
				tag: tag,
				children: []Node{
					TextNode(html.EscapeString(pp.Text)),
				},
			})
		case token.ListParagraph:
			currentIndentation := 1
			list := &HtmlNode{
				tag: "ul",
			}
			current := list

			for _, item := range pp.Items {
				if item.Indentation > currentIndentation {
					for i := currentIndentation; i < item.Indentation; i++ {
						child := &HtmlNode{
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
				current.AddChildren(&HtmlNode{
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
	safeStrings := map[string]template.HTML{}
	for k, v := range options.String {
		safeStrings[k] = template.HTML(v)
	}
	safeStrings["_Head"] = template.HTML(head.HTML(0))

	err = outTemplate.Execute(&out, struct {
		Body    template.HTML
		Bool    map[string]bool
		String  map[string]template.HTML
		Generic map[string]interface{}
	}{template.HTML(body.HTML(0)), options.Bool, safeStrings, options.Generic})
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

func fromLineContainer(line token.LineContainer) *HtmlNode {
	root := &HtmlNode{
		tag: "",
	}

	for _, text := range line.Tokens {
		switch tt := text.(type) {
		case token.TextToken:
			textRoot := &HtmlNode{
				tag: "",
			}
			current := textRoot

			if tt.Bold {
				node := &HtmlNode{
					tag: "b",
				}
				current.AddChildren(node)
				current = node
			}
			if tt.Italic {
				node := &HtmlNode{
					tag: "i",
				}
				current.AddChildren(node)
				current = node
			}
			if tt.Strike {
				node := &HtmlNode{
					tag: "s",
				}
				current.AddChildren(node)
				current = node
			}
			if tt.Link != "" {
				node := &HtmlNode{
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
			root.AddChildren(&HtmlNode{
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

func generateCode(p token.CodeParagraph, options eNote.Options) (*HtmlNode, string) {
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

	cssNode := &HtmlNode{
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
