package parser

import (
	"bufio"
	"eNote/token"
	"fmt"
	"io"
	"log"
)

var tokensMap = map[token.Type]token.Token{
	token.TypeBold:          token.BoldToken{},
	token.TypeItalic:        token.ItalicToken{},
	token.TypeLess:          token.LessToken{},
	token.TypeHeader:        token.HeaderToken{},
	token.TypeEqual:         token.EqualToken{},
	token.TypeSBracketOpen:  token.SBracketOpenToken{},
	token.TypeSBracketClose: token.SBracketCloseToken{},
	token.TypeTab:           token.TabToken{},
}

//Tokenizer parse a *os.File and return a slice of tokens
func Tokenizer(reader io.Reader) ([]token.Token, error) {
	tokenList := []token.Token{}
	r := bufio.NewReader(reader)
	buffer := ""
	currentEscape := false

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

		if n == token.TypeEscape && !currentEscape {
			currentEscape = true
			continue
		}

		if currentEscape && in(token.Type(n), token.WhitespaceEscape) {
			buffer += string(n)
			currentEscape = false
			continue
		}

		currentEscape = false

		if tok, ok := tokensMap[token.Type(n)]; ok {
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, tok)
		} else {
			switch n {
			case token.TypeNewLine:
				log.Println("\t- Found NewLineToken")
				if len(buffer) != 0 {
					addBufferToTokenBuffer(&tokenList, &buffer)
				}
				tokenList = append(tokenList, token.NewLineToken{})
			default:
				buffer += string(n)
			}
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

func in(t token.Type, slice []token.Type) bool {
	for _, tok := range slice {
		if t == tok {
			return true
		}
	}

	return false
}
