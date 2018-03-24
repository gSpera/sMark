package parser

import (
	"eNote/token"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

//ParseFile parse a *os.File
func ParseFile(fl *os.File) ([]token.Token, error) {
	tokenList := []token.Token{}
	char := make([]byte, 1)
	buffer := ""
	lineTokenBuffer := []token.Token{}
	currentLine := token.LineToken{}
	currentParagraph := token.ParagraphToken{}
	lastNewLine := false

	for {
		n, err := fl.Read(char)
		if n == 0 {
			addBufferToTokenBuffer(&lineTokenBuffer, &buffer)
			currentLine.Tokens = lineTokenBuffer
			currentParagraph.Lines = append(currentParagraph.Lines, currentLine)
			tokenList = append(tokenList, currentParagraph)
			fmt.Println("EOF")
			return tokenList, nil
		}
		if err != nil {
			return nil, err
		}

		switch char[0] {
		case '\t':
			if len(buffer) == 0 {
				currentLine.Indentation++
			}
		case '\n':
			if lastNewLine { //EOP: End Of Paragraph
				fmt.Println("EOP")
				tokenList = append(tokenList, currentParagraph)
				currentParagraph = token.ParagraphToken{}
			} else { //Generic newline
				fmt.Println("NewLine")
				// addBufferToTokenBuffer(&currentParagraph, &buffer)
				currentLine.Tokens = lineTokenBuffer
				spew.Dump(currentLine)
				currentParagraph.Lines = append(currentParagraph.Lines, currentLine)
				lineTokenBuffer = []token.Token{}
				currentLine = token.LineToken{}
			}

		case token.TypeBold:
			addBufferToTokenBuffer(&lineTokenBuffer, &buffer)
			fmt.Println("Bold")
			lineTokenBuffer = append(lineTokenBuffer, token.BoldToken{})
		case token.TypeItalic:
			addBufferToTokenBuffer(&lineTokenBuffer, &buffer)
			fmt.Println("Italic")
			lineTokenBuffer = append(lineTokenBuffer, token.ItalicToken{})

		default:
			// fmt.Printf("Char: %c\n", char[0])
			buffer += string(char[0])
		}
		lastNewLine = char[0] == '\n'
	}
}

func addBufferToTokenBuffer(tokenBuffer *[]token.Token, buffer *string) {
	*tokenBuffer = append(*tokenBuffer, token.TextToken{Text: *buffer})
	*buffer = ""
}

//ParseString parse a string
//NOT IMPLEMENTED
//TODO: Implement
func ParseString(str string) {}
