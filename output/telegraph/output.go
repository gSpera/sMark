package outTelegraph

import (
	"eNote/token"
	eNote "eNote/utils"
	"fmt"

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
		bold := false
		italic := false

		switch p.(type) {
		case token.HeaderParagraph:
			headerOptions := p.(token.HeaderParagraph).OptionsTemplate
			options.Update(headerOptions)
			if title == "" {
				title = *options.Title
			}
		case token.TitleParagraph:
			p := p.(token.TitleParagraph)

			if title == "" {
				title = p.Text.String()
			}

			fmt.Println("Title Paragraph")
			nodes = append(nodes, tgraph.NodeElement{
				Tag:      "h3",
				Children: []tgraph.Node{p.Text.String()},
			})

		case token.DivisorParagraph:
			fmt.Println("Divisor Paragraph")
			nodes = append(nodes, tgraph.NodeElement{
				Tag: "hr",
			})

		case token.TextParagraph:
			p := p.(token.TextParagraph)
			for _, line := range p.Lines {
				nodes = append(nodes, createLine(line, &bold, &italic)...)
				//Appending NewLine if the options allows it
				if *options.NewLine {
					nodes = append(nodes, tgraph.NodeElement{Tag: "br"})
				}
			}
			lastParagraph = p
		}

		if _, ok := lastParagraph.(token.TextParagraph); ok {
			nodes = append(nodes, tgraph.NodeElement{Tag: "br"})
		}

	}

	p := tgraph.Page{
		Title:   title,
		Content: nodes,
	}
	spew.Dump(p)
	return p
}

func createLine(line token.LineContainer, bold *bool, italic *bool) []tgraph.Node {
	nodes := []tgraph.Node{}
	for _, t := range line.Tokens {
		switch t.Type() {
		case token.TypeBold:
			*bold = !*bold
		case token.TypeItalic:
			*italic = !*italic
		case token.TypeText:
			switch {
			case *bold:
				nodes = append(nodes, tgraph.NodeElement{Tag: "b", Children: []tgraph.Node{t.String()}})
			case *italic:
				nodes = append(nodes, tgraph.NodeElement{Tag: "i", Children: []tgraph.Node{t.String()}})
			default:
				nodes = append(nodes, t.String())
			}
		}
	}

	fmt.Println("Create Line:")
	spew.Dump(nodes)
	return nodes
}
