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
			fmt.Println("EOF")
			return tokenList, nil
		}
		if err != nil {
			return nil, err
		}

		switch char[0] {
		case token.TypeBold:
			addBufferToTokenList(&tokenList, buffer)
			fmt.Println("Bold")
			tokenList = append(tokenList, token.BoldToken{})
		case token.TypeItalic:
			addBufferToTokenList(&tokenList, buffer)
			fmt.Println("Italic")
			tokenList = append(tokenList, token.ItalicToken{})

		default:
			// fmt.Printf("Char: %c\n", char[0])
			buffer += string(char[0])
		}
	}
}

func addBufferToTokenList(tokenList *[]token.Token, buffer string) {
	*tokenList = append(*tokenList, token.TextToken{Text: buffer})
}

//ParseString parse a string
//NOT IMPLEMENTED
//TODO: Implement
func ParseString(str string) {}
