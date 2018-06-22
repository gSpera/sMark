package token

import (
	"eNote/utils"
	"fmt"
)

//ParagraphToken is a special interface that indicates a token
type ParagraphToken interface{ IsParagraph() }

//HeaderParagraph is a pragraph containing header info
type HeaderParagraph struct {
	ParagraphToken
	eNote.Options
}

//TextParagraph rapresent a single paragraph
type TextParagraph struct {
	ParagraphToken
	Indentation int
	Lines       []LineContainer
}

//Type return the type of the Token
func (t TextParagraph) Type() Type {
	return TypeParagraph
}
func (t TextParagraph) String() string {
	return fmt.Sprintf("%s", t.Lines)
}

//DebugString is used for ast outputting
func (t TextParagraph) DebugString() string {
	return fmt.Sprintf("<TokenParagraph: %v>\n", func() string {
		str := ""
		for _, l := range t.Lines {
			str += l.DebugString()
		}
		return str
	}())
}

//DivisorParagraph is a pragraph containing a divisor line
type DivisorParagraph struct {
	ParagraphToken
}

//TitleParagraph is a paragraph rapresenting a Title
type TitleParagraph struct {
	ParagraphToken
	Text        string
	Indentation int
}

//SubtitleParagraph is a paragraph rapresenting a Subtitle,
//it is similar to TitleParagraph
type SubtitleParagraph struct {
	ParagraphToken
	Text        string
	Indentation int
}

//ListParagraph is a paragraph rapresenting a List
type ListParagraph struct {
	ParagraphToken
	Items []ListLine
}
