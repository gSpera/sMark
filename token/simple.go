package token

//SimpleToken is a generic interface for implementing simple tokens,
//a token is a rune which is divided by normal text
type SimpleToken interface {
	Token
	Char() rune
}

//IsToken is used internaly for rapresenting a token,
//calling this will probabaly panic

//BoldToken rapresent a / token
type BoldToken struct{ SimpleToken }

//ItalicToken rapresent a / token
type ItalicToken struct{ SimpleToken }

//HeaderToken rapresent a + token
type HeaderToken struct{ SimpleToken }

//NewLineToken rapresent a \n token
type NewLineToken struct{ SimpleToken }

//TabToken rapresent a \t token
type TabToken struct{ SimpleToken }

//LessToken rapresent a - token
type LessToken struct{ SimpleToken }

//EqualToken rapresent a - token
type EqualToken struct{ SimpleToken }

//SBracketOpenToken rapresent a [ token
type SBracketOpenToken struct{ SimpleToken }

//SBracketCloseToken rapresent a ] token
type SBracketCloseToken struct{ SimpleToken }

//EscapeToken is a special toekn used for escpaing other tokens, it rapresent a \ token
type EscapeToken struct{ SimpleToken }
