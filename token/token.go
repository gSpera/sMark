package token

import (
	"fmt"
)

//Type rappresent a type of token
type Type int8

//Token is an interface for token
type Token interface {
	fmt.Stringer

	//DebugString is used for ast outputting
	DebugString() string
	//Type is the type of the token
	Type() Type
}

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

	//TypeText is a text token
	TypeText = -1

	//TypeTokenLine rappresent a line
	TypeTokenLine
	//TypeParagraph rapresent a paragraph
	TypeParagraph
)

//UndefinedToken is a generic undefined token
type UndefinedToken struct{}

//Type returns the type of the Token
func (t UndefinedToken) Type() Type { return TypeUndefined }

func (t UndefinedToken) String() string { return "UNDEFINED" }

//DebugString is used for ast outputting
func (t UndefinedToken) DebugString() string { return "<UndefinedToken>" }

//BoldToken is a generic undefined token
type BoldToken struct{}

//Type returns the type of the Token
func (t BoldToken) Type() Type { return TypeBold }

func (t BoldToken) String() string { return "*" }

//DebugString is used for ast outputting
func (t BoldToken) DebugString() string { return "<BoldToken>" }

//ItalicToken is a generic undefined token
type ItalicToken struct{}

//Type returns the type of the Token
func (t ItalicToken) Type() Type { return TypeItalic }

func (t ItalicToken) String() string { return "/" }

//DebugString is used for ast outputting
func (t ItalicToken) DebugString() string { return "<ItalicToken>" }

//HeaderToken rapresent a char which compose the header line
type HeaderToken struct{}

//Type returns the type of the Token
func (t HeaderToken) Type() Type { return TypeHeader }

func (t HeaderToken) String() string { return "+" }

//DebugString is used for ast outputting
func (t HeaderToken) DebugString() string { return "<HeaderToken | +>" }

//NewLineToken is a generic undefined token
type NewLineToken struct{}

//Type returns the type of the Token
func (t NewLineToken) Type() Type { return TypeNewLine }

func (t NewLineToken) String() string { return "\n" }

//DebugString is used for ast outputting
func (t NewLineToken) DebugString() string { return "<NewLineToken>" }

//TabToken is a generic undefined token
type TabToken struct{}

//Type returns the type of the Token
func (t TabToken) Type() Type { return TypeTab }

func (t TabToken) String() string { return "\t" }

//DebugString is used for ast outputting
func (t TabToken) DebugString() string { return "<TabToken>" }

//TextToken a token conteining text
type TextToken struct {
	Text string
}

//Type returns the type of the Token
func (t TextToken) Type() Type { return TypeText }

func (t TextToken) String() string { return t.Text }

//DebugString is used for ast outputting
func (t TextToken) DebugString() string { return fmt.Sprintf("<TextToken: %s>\n", t.Text) }

//LessToken a token conteining a less(-) sign
type LessToken struct {
}

//Type returns the type of the Token
func (t LessToken) Type() Type { return TypeLess }

func (t LessToken) String() string { return string(TypeLess) }

//DebugString is used for ast outputting
func (t LessToken) DebugString() string { return fmt.Sprintf("<LessToken>") }

//EqualToken a token conteining a equal(=) sign
type EqualToken struct {
}

//Type returns the type of the Token
func (t EqualToken) Type() Type { return TypeEqual }

func (t EqualToken) String() string { return string(TypeEqual) }

//DebugString is used for ast outputting
func (t EqualToken) DebugString() string { return fmt.Sprintf("<EqualToken>") }

//LineState is the line metadata
type LineState struct {
	Indentation uint
}

//FromRune return an appropiate token for the rune
func FromRune(r rune) Token {
	switch r {
	case TypeBold:
		return &BoldToken{}
	case TypeItalic:
		return &ItalicToken{}
	default:
		return nil
	}
}
