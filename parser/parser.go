package parser

import (
	"io"
	"log"

	"github.com/gSpera/sMark/token"
	"github.com/gSpera/sMark/utils"
)

//ParseReader parse a source file and returns the paragraphs
func ParseReader(fl io.Reader) ([]token.ParagraphToken, error) {
	log.Println("Tokenizer")
	tokens, err := Tokenizer(fl)
	log.Println("Tokenizer DONE")
	if err != nil {
		return nil, err
	}

	log.Println("Searching Structures")
	tokens = TokenToStructure(tokens)
	log.Println("Searching Structures DONE")

	log.Println("Token To Line")
	lines := TokenToLine(tokens)
	log.Println("Token To Line DONE")

	log.Println("Line To Paragraph")
	paragraphs := TokenToParagraph(lines)
	log.Println("Line To Paragraph DONE")

	return paragraphs, nil
}

//OptionsFromParagraphs analyze the passed slice of paragraphs returning the final options contined in token.HeaderParagraphs.
//It return the optained options and a new list of paragraphs containing all the paragraphs but HeaderParagraphs
func OptionsFromParagraphs(paragraphs []token.ParagraphToken) (sMark.Options, []token.ParagraphToken) {
	options := sMark.NewOptions()
	newTokens := make([]token.ParagraphToken, 0, len(paragraphs))

	for _, p := range paragraphs {
		header, ok := p.(token.HeaderParagraph)
		if !ok {
			newTokens = append(newTokens, p)
			continue
		}

		options.Update(header.Options)
	}

	return options, newTokens
}

//TitleFromParagraph resolve a possible title from the passed paragraphs
func TitleFromParagraph(paragraph []token.ParagraphToken) string {
	titles := map[string]int{}

	for _, p := range paragraph {
		if pp, ok := p.(token.TitleParagraph); ok {
			txt := pp.Text
			value := pp.Indentation

			if currentValue, ok := titles[txt]; !ok || currentValue < value {
				titles[txt] = value
			}
		}

	}

	max := "No Title"
	maxV := 0

	for k, v := range titles {
		if v > maxV {
			max = k
			maxV = v
		}
	}

	return max
}
