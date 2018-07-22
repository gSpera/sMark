package token

import (
	"eNote/utils"
)

//ParagraphToken is a special interface that indicates a paragraph
type ParagraphToken interface {
	Token
	IsParagraph()
}

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

//CodeParagraph is a paragraph containing code
type CodeParagraph struct {
	ParagraphToken
	Lang string
	Text TextParagraph
}
