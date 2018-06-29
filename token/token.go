package token

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
}

//String creates a string with the content of the TextToken
//It respect attributes
func (t TextToken) String() string {
	var str string
	addAttr := func(str string) string {
		if t.Bold {
			return str + "*"
		}
		if t.Italic {
			return str + "/"
		}
		if t.Strike {
			return str + "-"
		}
		return str
	}

	str = addAttr("")
	str += t.Text
	str = addAttr(str)
	return str
}

//CheckBoxToken rapresent a checkbox, it is composed by [Char] (SBRacketOpenToken,TextToken,SBRacketCloseToken)
type CheckBoxToken struct {
	Token
	Char rune
}
