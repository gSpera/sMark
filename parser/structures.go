package parser

import (
	"eNote/token"
)

const maxTokenDistance = 255

//TokenToStructure checks the slice of tokens for multi-token tokens, like bold/italic text
func TokenToStructure(tokens []token.Token) []token.Token {
	toSkip := 0 //toSkip contains to token that the for loop need to skip

	for i, tok := range tokens {
		if toSkip > 0 {
			toSkip--
			continue
		}

		switch ttok := tok.(type) {
		case token.BoldToken:
			_ = ttok
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Bold: true} }
			tokens = checkRangeStruct(token.BoldToken{}, fn, tokens, i)
			return tokens
		case token.ItalicToken:
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Italic: true} }
			tokens = checkRangeStruct(token.ItalicToken{}, fn, tokens, i)
		}
	}

	return tokens
}

//checkRangeStruct searchs for a strcture like <token><Text></token>, for example with Bold *Text*
//but it could be used also for other tokens like italic
func checkRangeStruct(ending token.Token, generateToken func(string) token.TextToken, tokens []token.Token, start int) []token.Token {
	var buffer string

	//Starting from next token, stop after maxTokenDistance or when tokens finish
	for i := start + 1; i < len(tokens) || i-start > maxTokenDistance; i++ {

		//Found Ending Token
		if tokens[i] == ending {
			if i == start { //The two tokens are adjacent
				return tokens
			}

			newTokens := make([]token.Token, start) //create new buffer
			copy(newTokens, tokens)
			newTokens = append(newTokens, generateToken(buffer))
			return newTokens
		}

		switch tt := tokens[i].(type) {
		//Found text which we can embed in the token
		case token.TextToken:
			buffer += tt.Text

		//Another token is not accepted inside bold
		default:
			return tokens
		}
	}

	//No ending token found
	return tokens
}
