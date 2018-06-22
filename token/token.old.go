package token

//This file contains legacy code,
//With updates this file will keep shrinking until being eliminated
import (
	"fmt"
)

//Type rappresent a type of token
type Type int8

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

//DebugString is used for ast outputting
func (t TextToken) DebugString() string { return fmt.Sprintf("<TextToken: %s>\n", t.Text) }

//LineState is the line metadata
type LineState struct {
	Indentation int
}

//Type is deprecated
func (t BoldToken) Type() Type { return TypeBold }

//Type is deprecate
//Type is deprecated
func (t ItalicToken) Type() Type { return TypeItalic }

//Type is deprecated
func (t NewLineToken) Type() Type { return TypeNewLine }

//Type is deprecated
func (t TabToken) Type() Type { return TypeTab }

//Type is deprecated
func (t HeaderToken) Type() Type { return TypeHeader }

//Type is deprecated
func (t LessToken) Type() Type { return TypeLess }

//Type is deprecated
func (t EqualToken) Type() Type { return TypeEqual }

//Type is deprecated
func (t TextToken) Type() Type { return TypeText }

func (t BoldToken) String() string    { return "<BoldToken>" }
func (t ItalicToken) String() string  { return "<ItalicToken>" }
func (t HeaderToken) String() string  { return "+" }
func (t NewLineToken) String() string { return "<NewLineToken>" }
func (t TabToken) String() string     { return "\t" }
func (t LessToken) String() string    { return "<LessToken>" }
func (t EqualToken) String() string   { return "<EqualToken>" }

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
