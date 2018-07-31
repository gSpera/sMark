package token

import (
	"fmt"
	"strings"
)

//Token is an interface for simple and complex tokens
type Token interface {
	IsToken()

	//Type returns the Type of the token
	//It is used for recognizing the Token
	Type() Type
}

//TextToken is a complex token which contains text and some attribute.
type TextToken struct {
	Token

	Indentation int
	Text        string

	//Attributes
	Bold   bool
	Italic bool
	Strike bool
	Link   string
}

//String creates a string with the content of the TextToken
//It respect attributes
func (t TextToken) String() string {
	return t.string(func(s string) string { return s })
}

//StringEscape creates a string with the content of the TextToken
//It respect attributes, it escapes with backslash if it find escaped char
func (t TextToken) StringEscape() string {
	return t.string(EscapeString)
}

func (t TextToken) string(applyToText func(string) string) string {
	var str string
	addAttr := func(str string) string {
		var attr string
		if t.Bold {
			attr += "*"
		}
		if t.Italic {
			attr += "/"
		}
		if t.Strike {
			attr += "-"
		}
		if t.Link != "" {
			attr += "\""
		}
		return str + attr
	}

	str = strings.Repeat("\t", t.Indentation)
	str += addAttr("")
	str += applyToText(t.Text)
	str = addAttr(str)

	if t.Link != "" {
		str += fmt.Sprintf("@\"%s\"", t.Link)
	}

	return str
}

//CheckBoxToken rapresent a checkbox, it is composed by [Char] (SBRacketOpenToken,TextToken,SBRacketCloseToken)
type CheckBoxToken struct {
	Token
	Char rune
}
