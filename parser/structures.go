package parser

import (
	"eNote/token"
	"log"
)

const maxTokenDistance = 255

//TokenToStructure checks the slice of tokens for multi-token tokens, like bold/italic text
func TokenToStructure(tokens []token.Token) []token.Token {
	toSkip := 0 //toSkip contains to token that the for loop need to skip
	var newTokens []token.Token
	for i, tok := range tokens {
		if toSkip > 0 {
			toSkip--
			continue
		}

		switch tok.(type) {
		case token.BoldToken:
			log.Println("\t-Found BoldToken")
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Bold: true} }
			t, skip := checkRangeStruct(token.BoldToken{}, fn, tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}
			toSkip = skip
			newTokens = append(newTokens, t)
		case token.ItalicToken:
			log.Println("\t-Found ItalicToken")
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Italic: true} }
			t, skip := checkRangeStruct(token.ItalicToken{}, fn, tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}
			toSkip = skip
			newTokens = append(newTokens, t)
		default:
			newTokens = append(newTokens, tok)
		}
	}

	return newTokens
}

//checkRangeStruct searchs for a strcture like <token><Text></token>, for example with Bold *Text*
//but it could be used also for other tokens like italic
//returns the new token and tokens to skips, -1 if no structure found
func checkRangeStruct(ending token.Token, generateToken func(string) token.TextToken, tokens []token.Token, start int) (token.Token, int) {
	var buffer string

	//Starting from next token, stop after maxTokenDistance or when tokens finish
	for i := start + 1; i < len(tokens) || i-start > maxTokenDistance; i++ {

		//Found Ending Token
		if tokens[i] == ending {
			if i == start+1 { //The two tokens are adjacent
				return ending, -1
			}
			return generateToken(buffer), i - start
		}

		switch tt := tokens[i].(type) {
		//Found text which we can embed in the token
		case token.TextToken:
			buffer += tt.Text

		//Another token is not accepted inside
		default:
			return ending, -1
		}
	}

	//No ending token found
	return ending, -1
}
