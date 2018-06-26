package outTelegraph

import (
	"eNote/token"
	"eNote/utils"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	tgraph "github.com/toby3d/telegraph"
)

//ToString creates a telegraph that can be pubblished
//Options Meaning:
// - Only Body  has no meaning
// - CustomCSS  has no meaning
// - InlineCSS  has no meaning
// - EnableFont has no meaning
func ToString(paragraphs []token.ParagraphToken, options eNote.Options) tgraph.Page {
	nodes := []tgraph.Node{}
	title := *options.Title
	var lastParagraph token.ParagraphToken

	for _, p := range paragraphs {
		switch pp := p.(type) {
		case token.HeaderParagraph:
			panic("HeaderParagraph in output engine")
		case token.TitleParagraph:
			fmt.Println("Title Paragraph")
			nodes = append(nodes, createTitle(pp.Text))
		case token.SubtitleParagraph:
			fmt.Println("Subtitle Paragraph")
			nodes = append(nodes, createSubtitle(pp.Text))
		case token.DivisorParagraph:
			fmt.Println("Divisor Paragraph")
			nodes = append(nodes, createTag("hr"))

		case token.TextParagraph:
			fmt.Println("New Paragraph")
			p := p.(token.TextParagraph)
			par := createTag("p")
			for _, line := range p.Lines {
				fmt.Println("Appending Line")

				for _, tok := range line.Tokens {
					switch t := tok.(type) {
					case token.TextToken:
						par.Children = append(par.Children, createLine(t))
					default:
						if st, ok := t.(token.SimpleToken); ok {
							par.Children = append(par.Children, string(st.Char()))
							continue
						}
						panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", t, t))
					}
				}

				//Appending NewLine if the options allows it
				if *options.NewLine {
					fmt.Println("Adding NewLine")
					par.Children = append(par.Children, "\n")
				}
			}

			fmt.Println("Finished Paragraph")
			log.Println(spew.Sdump(par))
			nodes = append(nodes, par)
			lastParagraph = p
		case token.ListParagraph:
			fmt.Println("ListParagraph")
			list := createTag("ul")
			for _, item := range pp.Items {
				fmt.Println("Appending Item")
				li := createLi(item.Text.String())
				list.Children = append(list.Children, li)
			}

			spew.Dump(list)
			nodes = append(nodes, list)
			lastParagraph = pp
		}

		if _, ok := lastParagraph.(token.TextParagraph); ok {
			// nodes = append(nodes, tgraph.NodeElement{Tag: "br"})
		}

	}

	p := tgraph.Page{
		Title:   title,
		Content: nodes,
	}
	spew.Dump(p)
	return p
}

func createLine(text token.TextToken) tgraph.Node {
	fmt.Println("createLine")
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

	fmt.Println("Create Line:")
	return node
}

func createLi(text string) tgraph.Node {
	li := createTag("li")
	li.Children = []tgraph.Node{text}
	return li
}
