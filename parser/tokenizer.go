package parser

import (
	"eNote/token"
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
)

//Tokenizer parse a *os.File and return a slice of tokens
func Tokenizer(reader io.Reader) ([]token.Token, error) {
	tokenList := []token.Token{}
	char := make([]byte, 1)
	buffer := ""

	for {
		n, err := reader.Read(char)
		if n == 0 {
			addBufferToTokenBuffer(&tokenList, &buffer)
			fmt.Println("EOF")
			return tokenList, nil
		}
		if err != nil {
			return nil, err
		}

		switch char[0] {
		case '\t':
			fmt.Println("TAB")
			tokenList = append(tokenList, token.TabToken{})
		case '\n':
			fmt.Println("NewLine")
			if len(buffer) != 0 {
				addBufferToTokenBuffer(&tokenList, &buffer)
			}
			tokenList = append(tokenList, token.NewLineToken{})
		case token.TypeBold:
			fmt.Println("Bold")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.BoldToken{})

		case token.TypeItalic:
			fmt.Println("Italic")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.ItalicToken{})

		default:
			// fmt.Printf("Char: %c\n", char[0])
			buffer += string(char[0])
		}
	}
}

func addBufferToTokenBuffer(tokenBuffer *[]token.Token, buffer *string) {
	*tokenBuffer = append(*tokenBuffer, token.TextToken{Text: *buffer})
	*buffer = ""
}

//TokenToLine divide a slice of tokens in lines
func TokenToLine(tokens []token.Token) []token.LineToken {
	lines := []token.LineToken{}
	currentLine := token.LineToken{}
	indent := true

	for _, t := range tokens {
		switch t.(type) {
		case token.TabToken:
			if indent {
				currentLine.Indentation++
			} else {
				indent = false
			}
		case token.NewLineToken:
			fmt.Println("NewLine")
			spew.Dump(currentLine.Tokens)
			lines = append(lines, currentLine)
			currentLine = token.LineToken{}
		default:
			currentLine.Tokens = append(currentLine.Tokens, t)
		}
	}

	return lines
}

//TokenToParagraph divide a slice of lines in paragraphs
func TokenToParagraph(lines []token.LineToken) []token.ParagraphToken {
	fmt.Printf("Paragraphs: %d Lines\n", len(lines))
	paragraphs := []token.ParagraphToken{}
	currentParagraph := token.ParagraphToken{}

	for i, t := range lines {
		if i == 0 {
			continue
		}
		previousToken := lines[i-1]

		currentEmpty := len(t.Tokens) == 0
		previuosEmpty := len(previousToken.Tokens) == 0
		fmt.Println(len(t.Tokens), len(previousToken.Tokens))
		spew.Dump(t.Tokens)

		if currentEmpty && previuosEmpty {
			fmt.Println("New Paragraph")
			paragraphs = append(paragraphs, currentParagraph)
			currentParagraph = token.ParagraphToken{}
		} else {
			currentParagraph.Lines = append(currentParagraph.Lines, t)
		}
	}

	paragraphs = append(paragraphs, currentParagraph)
	return paragraphs
}
