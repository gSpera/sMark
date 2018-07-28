package outTelegraph

import (
	"eNote/token"
	"eNote/utils"
	"fmt"
	"log"
	"strings"

	tgraph "github.com/toby3d/telegraph"
)

//ToString creates a telegraph that can be pubblished
func ToString(paragraphs []token.ParagraphToken, options eNote.Options) tgraph.Page {
	nodes := []tgraph.Node{}
	title := options.String["Title"]

	for _, p := range paragraphs {
		switch pp := p.(type) {
		case token.HeaderParagraph:
			panic("HeaderParagraph in output engine")
		case token.TitleParagraph:
			log.Println("\t- Title Paragraph")
			nodes = append(nodes, createTitle(pp.Text))
		case token.SubtitleParagraph:
			log.Println("\t- Subtitle Paragraph")
			nodes = append(nodes, createSubtitle(pp.Text))
		case token.DivisorParagraph:
			log.Println("\t- Divisor Paragraph")
			nodes = append(nodes, createTag("hr"))

		case token.TextParagraph:
			log.Println("\t- New Paragraph")
			p := p.(token.TextParagraph)
			par := createTag("p")
			var quote bool
			quoteBlock := createTag("blockquote")

			for _, line := range p.Lines {
				if !line.Quote && quote { //End of block of quote
					par.Children = append(par.Children, quoteBlock)
					quoteBlock = createTag("blockquote")
				}
				quote = line.Quote
				if quote {
					quoteBlock.Children = append(quoteBlock.Children, fromLineContainer(line).Children...)
				} else {
					par.Children = append(par.Children, fromLineContainer(line).Children...)
				}

				//Appending NewLine if the options allows it
				if options.Bool["NewLine"] && len(par.Children) != 0 {
					node := &par
					if quote {
						node = &quoteBlock
					}
					node.Children = append(node.Children, "\n")
				}
			}

			log.Println("\t- Finished Paragraph")
			if len(par.Children) != 0 {
				nodes = append(nodes, par)
			}
		case token.ListParagraph:
			log.Println("\t- ListParagraph")
			tre := &tree{createTag("ul"), nil}
			root := tre
			currentIndentation := 0

			//Telegra.ph doesn't support nested ul so this is pretty useless,
			//maybe in a future update they will add them
			for _, item := range pp.Items {
				if item.Indentation > currentIndentation {
					for i := currentIndentation; i < item.Indentation; i++ {
						tr := &tree{createTag("ul"), tre}
						tre.Children = append(tre.Children, tr)
						tre = tr
					}
				} else if item.Indentation < currentIndentation {
					for i := currentIndentation; i > item.Indentation; i-- {
						tre = tre.parent
					}
				}

				li := createTag("li")
				//Emulate nested list
				if item.Indentation > 1 {
					stars := createBold(strings.Repeat("*", item.Indentation-1))
					li.Children = append(li.Children, stars)
					li.Children = append(li.Children, ": ")
				}
				li.Children = append(li.Children, fromLineContainer(item.Text))
				tre.Children = append(tre.Children, li)
			}

			nodes = append(nodes, root)
		case token.CodeParagraph:
			log.Println("\t- CodeParagraph")
			pre := createTag("pre")
			for _, line := range pp.Text.Lines {
				pre.Children = append(pre.Children, line.String()+"\n")
			}
			nodes = append(nodes, pre)
		}
	}

	p := tgraph.Page{
		Title:   title,
		Content: nodes,
	}
	return p
}

func createLine(text token.TextToken) tgraph.Node {
	var node tgraph.Node = text.Text
	if text.Bold {
		node = tgraph.NodeElement{Tag: "b", Children: []tgraph.Node{node}}
	}
	if text.Italic {
		node = tgraph.NodeElement{Tag: "b", Children: []tgraph.Node{node}}
	}
	if text.Strike {
		node = tgraph.NodeElement{Tag: "s", Children: []tgraph.Node{node}}
	}

	return node
}

func createLi(text tgraph.Node) tgraph.Node {
	li := createTag("li")
	li.Children = []tgraph.Node{text}
	return li
}

func createCheckBox(t token.CheckBoxToken) tgraph.Node {
	code := createTag("code")
	code.Children = []tgraph.Node{fmt.Sprintf("[%c]", t.Char)}
	return code
}

func fromLineContainer(line token.LineContainer) tgraph.NodeElement {
	res := createTag("p")

	for _, tok := range line.Tokens {
		switch t := tok.(type) {
		case token.TextToken:
			res.Children = append(res.Children, createLine(t))
		case token.SimpleToken:
			res.Children = append(res.Children, string(t.Char()))
		case token.CheckBoxToken:
			res.Children = append(res.Children, createCheckBox(t))
		default:
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", t, t))
		}
	}

	return res
}
