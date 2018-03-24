package token

import (
	"fmt"
)

//Type rappresent a type of token
type Type int8

//Token is an interface for token
type Token interface {
	fmt.Stringer

	//The type of the token
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
)

//UndefinedToken is a generic undefined token
type UndefinedToken struct{}

//Type returns the type of the Token
func (t UndefinedToken) Type() Type { return TypeUndefined }

//BoldToken is a generic undefined token
type BoldToken struct{}

//Type returns the type of the Token
func (t BoldToken) Type() Type { return TypeBold }

func (t BoldToken) String() string { return "<BoldToken>\n" }

//ItalicToken is a generic undefined token
type ItalicToken struct{}

//Type returns the type of the Token
func (t ItalicToken) Type() Type { return TypeItalic }

func (t ItalicToken) String() string { return "<ItalicToken>\n" }

//TextToken a token conteining text
type TextToken struct {
	Text string
}

//Type returns the type of the Token
func (t TextToken) Type() Type { return TypeText }

func (t TextToken) String() string { return fmt.Sprintf("<TextToken: %s>\n", t.Text) }

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
