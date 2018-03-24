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

	for {
		n, err := fl.Read(char)
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
			addBufferToTokenBuffer(&tokenList, &buffer)
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

//ParseString parse a string
//NOT IMPLEMENTED
//TODO: Implement
func ParseString(str string) {}
