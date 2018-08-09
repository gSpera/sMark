package parser

import (
	"eNote/token"
	"eNote/utils"
	"io"
	"log"
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

//parseHeader parses a header string returning the Key and Value
func parseHeader(line string) (string, string) {
	key := ""
	buffer := []rune{}

	for i, ch := range line {
		switch ch {
		case '=':
			key = string(buffer)
			buffer = []rune{}
		case ';': //Comment
			if line[i-1] == ' ' {
				return key, string(buffer)
			}
			buffer = append(buffer, ch)
		case ' ':
			//Ignore if last character
			if len(line) <= i+1 {
				continue
			}

			if line[i-1] == '=' || line[i+1] == '=' { //Space are allowed around equal sign
				continue
			}

			if len(line) > i && line[i+1] == ';' { //Comment incoming
				continue
			}

			//Else insert it
			fallthrough
		default:
			buffer = append(buffer, ch)
		}
	}

	return key, string(buffer)
}

//OptionsFromParagraphs analyze the passed slice of paragraphs returning the final options contined in token.HeaderParagraphs.
//It return the optained options and a new list of paragraphs containing all the paragraphs but HeaderParagraphs
func OptionsFromParagraphs(paragraphs []token.ParagraphToken) (eNote.Options, []token.ParagraphToken) {
	options := eNote.NewOptions()
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
