package token

//This file contains legacy code,
//With updates this file will keep shrinking until being eliminated

//Type rappresent a type of token
type Type int

//Types definition
//They can rapresent the rapresented char but it is not mandatory
const (
	//TypeUndefined is an undefined token
	TypeUndefined = iota
	//TypeBold is an bold starting/ending token
	TypeBold = '*'
	//TypeItalic is an italic starting/ending token
	TypeItalic = '/'

	//TypeNewLine is a newline
	TypeNewLine = '\n'
	//TypeTab is a tab
	TypeTab = '\t'

	//TypeHeader is the token used in an header
	TypeHeader = '+'

	//TypeLess is the token used for subtitle and for divisor
	TypeLess = '-'

	//TypeEqual is the token used for title
	TypeEqual = '='

	TypeSBracketOpen  = '['
	TypeSBracketClose = ']'
	TypeEscape        = '\\'

	//TypeText is a text token
	TypeText = -1

	TypePipe  = '|'
	TypeQuote = '"'
	TypeAt    = '@'

	TypeParagraphHeader
	TypeParagraphText
	TypeParagraphDivisor
	TypeParagraphTitle
	TypeParagraphSubtitle
	TypeParagraphList
	TypeCheckBox
)

//WhitespaceEscape is a slice of types that can be escaped from the EscapeToken
var WhitespaceEscape = []Type{
	TypeBold,
	TypeItalic,
	TypeTab,
	TypeLess,
}

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

//Type returns the type of the Token
func (p SBracketOpenToken) Type() Type { return TypeSBracketOpen }

//Type returns the type of the Token
func (p SBracketCloseToken) Type() Type { return TypeSBracketClose }

//Type returns the type of the Token
func (p CheckBoxToken) Type() Type { return TypeCheckBox }

//Type returns the type of the Token
func (p EscapeToken) Type() Type { return TypeEscape }

//Type returns the type of the Token
func (p PipeToken) Type() Type { return TypePipe }

//Type returns the type of the Token
func (p QuoteToken) Type() Type { return TypeQuote }

//Type returns the type of the Token
func (p AtToken) Type() Type { return TypeAt }
