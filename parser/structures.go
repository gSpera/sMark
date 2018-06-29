package parser

import (
	"eNote/token"
	"log"
	"unicode/utf8"
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
		case token.LessToken:
			log.Println("\t-Found LessToken")
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Strike: true} }
			t, skip := checkRangeStruct(token.LessToken{}, fn, tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}
			toSkip = skip
			newTokens = append(newTokens, t)
		case token.SBracketOpenToken:
			log.Println("\t-Found Opening Square Bracket")
			t, skip := searchCheckbox(tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}
			log.Println("\t-Found CheckBox")
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

func searchCheckbox(tokens []token.Token, start int) (token.CheckBoxToken, int) {
	var res token.CheckBoxToken

	//Check opening bracket
	if tokens[start].Type() != token.TypeSBracketOpen {
		log.Println("\t\t- Start is not a SBracketOpenToken")
		return res, -1
	}

	//Check closing bracket
	if tokens[start+2].Type() != token.TypeSBracketClose {
		log.Println("\t\t- Start+1 is not a SBracketCloseToken")
		return res, -1
	}
	if _, ok := tokens[start+1].(token.TextToken); !ok {
		log.Println("\t\t- Start is not a String size is not 1")
		return res, -1
	}

	text := tokens[start+1].(token.TextToken)

	if utf8.RuneCountInString(text.Text) != 1 { //Text inside Brackets must be only one char
		return res, -1
	}

	char, _ := utf8.DecodeRuneInString(text.Text)
	return token.CheckBoxToken{Char: char}, 2
}
