package parser

import (
	"bufio"
	"eNote/token"
	"fmt"
	"io"
	"log"
)

//Tokenizer parse a *os.File and return a slice of tokens
func Tokenizer(reader io.Reader) ([]token.Token, error) {
	tokenList := []token.Token{}
	r := bufio.NewReader(reader)
	buffer := ""

	for {
		n, size, err := r.ReadRune()

		if size == 0 {
			addBufferToTokenBuffer(&tokenList, &buffer)
			//Adds new line at end if no one
			if tokenList[len(tokenList)-1].Type() != token.TypeNewLine {
				tokenList = append(tokenList, token.NewLineToken{})
			}
			fmt.Println("EOF")
			tokenList = checkTokenList(tokenList)
			return tokenList, nil
		}
		if err != nil {
			return nil, err
		}

		switch n {
		case token.TypeTab:
			log.Println("\t- Found TabToken")
			tokenList = append(tokenList, token.TabToken{})
		case token.TypeNewLine:
			log.Println("\t- Found NewLineToken")
			if len(buffer) != 0 {
				addBufferToTokenBuffer(&tokenList, &buffer)
			}
			tokenList = append(tokenList, token.NewLineToken{})
		case token.TypeBold:
			log.Println("\t- Found BoldToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.BoldToken{})
		case token.TypeItalic:
			log.Println("\t- Found ItalicToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.ItalicToken{})
		case token.TypeLess:
			log.Println("\t- Found LessToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.LessToken{})
		case token.TypeHeader:
			log.Println("\t- Found HeaderToken")
			tokenList = append(tokenList, token.HeaderToken{})
		case token.TypeEqual:
			log.Println("\t- Found EqualToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.EqualToken{})
		default:
			buffer += string(n)
		}
	}
}

func addBufferToTokenBuffer(tokenBuffer *[]token.Token, buffer *string) {
	if len(*buffer) == 0 {
		return
	}

	*tokenBuffer = append(*tokenBuffer, token.TextToken{Text: *buffer})
	*buffer = ""
}

func checkTokenList(tokenList []token.Token) []token.Token {
	if len(tokenList) == 0 {
		return tokenList
	}

	//Remove first TextToken if it is empty
	if t, ok := tokenList[0].(token.TextToken); ok && len(t.Text) == 0 {
		return tokenList[1:]
	}

	return tokenList
}
