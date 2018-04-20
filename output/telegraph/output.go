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
			fmt.Println("New Paragraph")
			p := p.(token.TextParagraph)
			par := tgraph.NodeElement{Tag: "p"}
			var fullLine []token.Token
			for _, line := range p.Lines {
				fmt.Println("Appending Line")
				fullLine = append(fullLine, line.Tokens...)

				//Appending NewLine if the options allows it
				if *options.NewLine {
					fmt.Println("Adding NewLine")
					fullLine = append(fullLine, token.TextToken{Text: "\n"})
				}
			}

			par.Children = append(par.Children, createLine(fullLine, &bold, &italic)...)
			fmt.Println("Finished Paragraph")
			spew.Dump(par)
			nodes = append(nodes, par)
			lastParagraph = p
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

func createLine(line []token.Token, bold *bool, italic *bool) []tgraph.Node {
	nodes := []tgraph.Node{}
	fmt.Println("createLine")
	spew.Dump(line)
	for _, t := range line {
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
