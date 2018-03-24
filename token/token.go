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
	//TypeText is a text token
	TypeText

	//TypeTokenLine rappresent a line
	TypeTokenLine
)

//UndefinedToken is a generic undefined token
type UndefinedToken struct{}

//Type returns the type of the Token
func (t UndefinedToken) Type() Type { return TypeUndefined }

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

//LineToken is a token which rappresent a list of Tokens with some attributes
type LineToken struct {
	LineState
	Tokens []Token
}

//Type return the type of the Token
func (t LineToken) Type() Type { return TypeTokenLine }

func (t LineToken) String() string {
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
func (t LineToken) DebugString() string {
	return fmt.Sprintf("<LineToken: {%v}%s>\n", t.LineState, t.Tokens)
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
