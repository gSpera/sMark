package parser

import (
	"eNote/token"
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

//TokenToParagraph divide a slice of lines in paragraphs
func TokenToParagraph(lines []token.LineToken) []token.ParagraphToken {
	fmt.Printf("Paragraphs: %d Lines\n", len(lines))
	paragraphs := []token.ParagraphToken{}
	currentParagraph := token.TextParagraph{}
	header := false
	var list []token.ListLine
	var skip int

	for i, t := range lines {
		if skip > 0 {
			skip--
			continue
		}
		var lastLine token.LineToken
		if i == 0 {
			lastLine = nil
		} else {
			lastLine = lines[i-1]
		}

		if _, ok := t.(token.ListLine); !ok && list != nil {
			paragraphs = append(paragraphs, token.ListParagraph{Items: list})
			list = nil
		}

		switch tt := t.(type) {
		case token.EqualLine:
			log.Println("\t- Found EqualLine")
			fmt.Println(len(currentParagraph.Lines))
			p, ok := equalLine(currentParagraph, lines, i)
			if ok {
				paragraphs = append(paragraphs, p)
				currentParagraph = token.TextParagraph{}
			} else {
				currentParagraph.Lines = append(currentParagraph.Lines,
					simpleText(strings.Repeat("=", tt.Length)))
			}
		case token.HeaderLine:
			fmt.Println("HeaderLine")
			if !header && len(currentParagraph.Lines) != 0 {
				fmt.Println("Not at first")
				continue
			}

			if header {
				log.Println("\t- Found Header Paragraph")
				paragraphs = append(paragraphs, token.HeaderParagraph{Options: parseHeaderLines(currentParagraph)})
				currentParagraph = token.TextParagraph{}
			}

			header = !header

		case token.LessLine:
			log.Println("\t- Found LessLine")
			fmt.Println(len(currentParagraph.Lines))
			p, ok := lessLine(currentParagraph, lines, i)
			if ok {
				log.Println("\t\t- Subtitle")
				paragraphs = append(paragraphs, p)
				currentParagraph = token.TextParagraph{}
			} else {
				log.Println("\t\t- Not Subtitle")
				currentParagraph.Lines = append(currentParagraph.Lines,
					simpleText(strings.Repeat("-", tt.Length)))
			}

		case token.ListLine:
			if len(currentParagraph.Lines) != 0 {
				paragraphs = append(paragraphs, currentParagraph)
			}

			list := searchList(lines, i)
			skip = len(list.Items)
			paragraphs = append(paragraphs, list)
			currentParagraph = token.TextParagraph{}

		case token.LineContainer:
			fmt.Println("LineContainer")
			currentEmpty := len(tt.Tokens) == 0
			spew.Dump(tt.Tokens)

			//Check if line is empty causing End Of Paragraph
			if l, ok := lastLine.(token.LineContainer); currentEmpty && ok && len(l.Tokens) != 0 {
				log.Println("\t- Found Text Paragraph")
				fmt.Printf("LastLine: Len: %d %T: %+v\n", len(l.Tokens), l, l)
				spew.Dump(l.Tokens)
				checkIndentation(&currentParagraph)
				fmt.Printf("Indentation after: %d\n", currentParagraph.Indentation)
				paragraphs = append(paragraphs, currentParagraph)
				currentParagraph = token.TextParagraph{}
			} else {
				fmt.Println("Old Paragraph")
				spew.Dump(tt.Tokens)
				currentParagraph.Lines = append(currentParagraph.Lines, tt)
			}
		default:
			panic(fmt.Sprintf("Line=>Paragraph for %T{%+v} not defined", tt, tt))
		}

	}

	checkIndentation(&currentParagraph)
	paragraphs = append(paragraphs, currentParagraph)
	return paragraphs
}

//notEmptyLines returns the number of lines which are not empty
func notEmptyLines(lines []token.LineContainer) uint {
	notEmpty := uint(0)

	for _, l := range lines {
		if len(l.Tokens) != 0 {
			notEmpty++
		}
	}

	return notEmpty
}

//equalLine searches for paragraph made with equalLine
//returns if found a new paragraph or not
func equalLine(currentParagraph token.TextParagraph, lines []token.LineToken, index int) (token.ParagraphToken, bool) {
	currentLine := lines[index].(token.EqualLine)
	//Need a line over
	if index == 0 {
		return nil, false
	}

	//Seartch if line number is right
	if notEmptyLines(currentParagraph.Lines) != 1 {
		log.Println("\t\t- Wrong number of lines:", notEmptyLines(currentParagraph.Lines))
		return nil, false
	}

	if _, ok := lines[index-1].(token.LineContainer); !ok {
		log.Println("\t\t- LastLine is not a token.LineContainer")
		return nil, false
	}

	lastLine := lines[index-1].(token.LineContainer)

	if len(lastLine.Tokens) == 0 {
		log.Println("\t\t- LastLine is empty")
		return nil, false
	}

	if lastLine.Indentation != currentLine.Indentation {
		log.Println("\t\t- Indentation are differents", lastLine.Indentation, lastLine.Indentation)
		return nil, false
	}

	log.Println("\t- Found Title Paragraph")
	return token.TitleParagraph{
		Text:        lastLine.StringNoTab(),
		Indentation: lastLine.Indentation,
	}, true
}

//lessLine searches for paragraph made with lessLine
//returns if found a new paragraph or not
func lessLine(currentParagraph token.TextParagraph, lines []token.LineToken, index int) (token.ParagraphToken, bool) {
	//Need a line over
	if index == 0 {
		log.Println("\t\t\t- Index is 0")
		return nil, false
	}

	if _, ok := lines[index-1].(token.LineContainer); !ok {
		log.Println("\t\t\t- Lne before is not TextToken")
		return nil, false
	}
	lastLine := lines[index-1].(token.LineContainer)

	if len(lastLine.Tokens) == 0 {
		log.Println("\t- Found Divisor")
		return token.DivisorParagraph{}, true
	}

	if notEmptyLines(currentParagraph.Lines) != 1 {
		// log.Println(currentParagraph)
		log.Println("\t\t- Wrong number of lines:", notEmptyLines(currentParagraph.Lines))
		return nil, false
	}

	if len(lastLine.Tokens) == 0 {
		log.Println("\t\t- LastLine is empty")
		return nil, false
	}

	if lastLine.Indentation != lastLine.Indentation {
		log.Println("\t\t- Indentation are differents", lastLine.Indentation, lastLine.Indentation)
		return nil, false
	}

	log.Println("\t- Found Subtitle Paragraph")
	return token.SubtitleParagraph{
		Text:        lastLine.StringNoTab(),
		Indentation: lastLine.Indentation,
	}, true
}

//searchList generate a list paragraph
func searchList(lines []token.LineToken, index int) token.ListParagraph {
	list := token.ListParagraph{}

	for i := index; ; i++ {
		if i >= len(lines) { //EOF
			break
		}

		if line, ok := lines[i].(token.ListLine); ok {
			list.Items = append(list.Items, line)
		} else {
			break
		}
	}

	return list
}

func simpleText(text string) token.LineContainer {
	return token.LineContainer{
		Tokens: []token.Token{
			token.TextToken{Text: text},
		},
	}
}