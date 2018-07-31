package token

import "strings"

//Char returns *
func (t BoldToken) Char() rune { return '*' }

//Char returns /
func (t ItalicToken) Char() rune { return '/' }

//Char returns +
func (t HeaderToken) Char() rune { return '+' }

//Char returns \n
func (t NewLineToken) Char() rune { return '\n' }

//Char returns \t
func (t TabToken) Char() rune { return '\t' }

//Char returns -
func (t LessToken) Char() rune { return '-' }

//Char returns =
func (t EqualToken) Char() rune { return '=' }

//Char returns [
func (t SBracketOpenToken) Char() rune { return '[' }

//Char returns ]
func (t SBracketCloseToken) Char() rune { return ']' }

//Char returns |
func (t PipeToken) Char() rune { return '|' }

//Char returns "
func (t QuoteToken) Char() rune { return '"' }

//Char returns @
func (t AtToken) Char() rune { return '@' }

//EscapeString escapes the current string from charachetrs like * or /
func EscapeString(str string) string {
	for _, ch := range WhitespaceEscape {
		str = strings.Replace(str, string(ch), "\\"+string(ch), -1)
	}

	return str
}
