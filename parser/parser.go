package parser

import (
	"eNote/token"
	"fmt"
	"os"
)

//ParseFile parse a *os.File
func ParseFile(fl *os.File) ([]token.Token, error) {
	tokenList := []token.Token{}
	char := make([]byte, 1)
	buffer := ""
	tokenBuffer := []token.Token{}

	currentLine := token.LineState{}
	for {
		n, err := fl.Read(char)
		if n == 0 {
			addBufferToTokenBuffer(&tokenBuffer, &buffer)
			tokenList = append(tokenList, token.LineToken{LineState: currentLine, Tokens: tokenBuffer})
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
			addBufferToTokenBuffer(&tokenBuffer, &buffer)
			tokenList = append(tokenList, token.LineToken{LineState: currentLine, Tokens: tokenBuffer})
			tokenBuffer = []token.Token{}
			currentLine = token.LineState{}

		case token.TypeBold:
			addBufferToTokenBuffer(&tokenBuffer, &buffer)
			fmt.Println("Bold")
			tokenBuffer = append(tokenBuffer, token.BoldToken{})
		case token.TypeItalic:
			addBufferToTokenBuffer(&tokenBuffer, &buffer)
			fmt.Println("Italic")
			tokenBuffer = append(tokenBuffer, token.ItalicToken{})

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

//ParseString parse a string
//NOT IMPLEMENTED
//TODO: Implement
func ParseString(str string) {}
