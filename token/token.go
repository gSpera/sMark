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

	//TypeText is a text token
	TypeText

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

//LineToken is a special interface that indicates a token
type LineToken interface{ IsToken() }

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

//LineState is the line metadata
type LineState struct {
	Indentation uint
}

//LineContainer is a token which rappresent a list of Tokens with some attributes
type LineContainer struct {
	LineState
	Tokens []Token
}

//Type return the type of the Token
func (t LineContainer) Type() Type { return TypeTokenLine }

func (t LineContainer) String() string {
	str := ""
	for i := uint(0); i < t.LineState.Indentation; i++ {
		str += "\t"
	}
	for _, t := range t.Tokens {
		str += fmt.Sprintf("%v", t)
	}
	return str + "\n"
}

//DebugString is used for ast outputting
func (t LineContainer) DebugString() string {
	return fmt.Sprintf("<LineToken: {%+v}%s>\n", t.LineState, t.Tokens)
}

//ParagraphToken rapresent a single paragraph
type ParagraphToken struct {
	Indentation int
	Lines       []LineContainer
}

//Type return the type of the Token
func (t ParagraphToken) Type() Type {
	return TypeParagraph
}
func (t ParagraphToken) String() string {
	return fmt.Sprintf("%s", t.Lines)
}

//DebugString is used for ast outputting
func (t ParagraphToken) DebugString() string {
	return fmt.Sprintf("<TokenParagraph: %v>\n", func() string {
		str := ""
		for _, l := range t.Lines {
			str += l.DebugString()
		}
		return str
	}())
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
