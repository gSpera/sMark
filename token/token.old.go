package token

//This file contains legacy code,
//With updates this file will keep shrinking until being eliminated

//Type rappresent a type of token
type Type int

//Default Types
const (
	TypeUndefined = iota
	TypeBold      = '*'
	TypeItalic    = '/'
	TypeNewLine   = '\n'
	TypeTab       = '\t'
	TypeHeader    = '+'
	TypeLess      = '-'
	TypeEqual     = '='
	TypeText      = -1

	TypeParagraphHeader
	TypeParagraphText
	TypeParagraphDivisor
	TypeParagraphTitle
	TypeParagraphSubtitle
	TypeParagraphList
)

//Defaults Type Method

//Type returns the type of the Token
func (t BoldToken) Type() Type { return TypeBold }

//Type returns the type of the Token
func (t ItalicToken) Type() Type { return TypeItalic }

//Type returns the type of the Token
func (t NewLineToken) Type() Type { return TypeNewLine }

//Type returns the type of the Token
func (t TabToken) Type() Type { return TypeTab }

//Type returns the type of the Token
func (t HeaderToken) Type() Type { return TypeHeader }

//Type returns the type of the Token
func (t LessToken) Type() Type { return TypeLess }

//Type returns the type of the Token
func (t EqualToken) Type() Type { return TypeEqual }

//Type returns the type of the Token
func (t TextToken) Type() Type { return TypeText }

//Type returns the type of the Token
func (p HeaderParagraph) Type() Type { return TypeParagraphHeader }

//Type returns the type of the Token
func (p TextParagraph) Type() Type { return TypeParagraphText }

//Type returns the type of the Token
func (p DivisorParagraph) Type() Type { return TypeParagraphDivisor }

//Type returns the type of the Token
func (p TitleParagraph) Type() Type { return TypeParagraphTitle }

//Type returns the type of the Token
func (p SubtitleParagraph) Type() Type { return TypeParagraphSubtitle }

//Type returns the type of the Token
func (p ListParagraph) Type() Type { return TypeParagraphList }

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
		return str
	}

	str = addAttr("")
	str += t.Text
	str = addAttr(str)
	return str
}
