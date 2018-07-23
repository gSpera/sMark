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
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Bold: true} }
			t, skip := checkRangeStruct(token.BoldToken{}, fn, tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}

			log.Println("\t-Found Bold Text")
			toSkip = skip
			newTokens = append(newTokens, t)
		case token.ItalicToken:
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Italic: true} }
			t, skip := checkRangeStruct(token.ItalicToken{}, fn, tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}

			log.Println("\t-Found Italic Text")
			toSkip = skip
			newTokens = append(newTokens, t)
		case token.LessToken:
			fn := func(buffer string) token.TextToken { return token.TextToken{Text: buffer, Strike: true} }
			t, skip := checkRangeStruct(token.LessToken{}, fn, tokens, i)
			if skip == -1 {
				newTokens = append(newTokens, tok)
				continue
			}

			log.Println("\t-Found Strike-Throught Text")
			toSkip = skip
			newTokens = append(newTokens, t)
		case token.SBracketOpenToken:
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

	if len(tokens)-start < 3 {
		log.Printf("Not enough tokens: %d (len(tokens)[%d] - start[%d] [= %d] < 3)\n", len(tokens), len(tokens), start, len(tokens)-start)
		return res, -1
	}

	//Check opening bracket
	if tokens[start].Type() != token.TypeSBracketOpen {
		panic("SBracketOpenToken not found")
	}

	//Check closing bracket
	if tokens[start+2].Type() != token.TypeSBracketClose {
		log.Println("\t\t- SBracketCloseToken not found")
		return res, -1
	}

	var char rune
	textToken, ok0 := tokens[start+1].(token.TextToken)
	simpleToken, ok1 := tokens[start+1].(token.SimpleToken)

	if ok0 {
		text := textToken.Text
		if utf8.RuneCountInString(text) != 1 { //Text inside Brackets must be only one char
			return res, -1
		}
		char, _ = utf8.DecodeRuneInString(text)
	} else if ok1 {
		char = simpleToken.Char()
	} else {
		log.Println("\t\t- TextToken not found")
		return res, -1
	}

	return token.CheckBoxToken{Char: char}, 2
}
