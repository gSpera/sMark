package token

import "fmt"

//ParagraphToken is a special interface that indicates a token
type ParagraphToken interface{ IsParagraph() }

//HeaderParagraph is a pragraph containing header info
type HeaderParagraph struct {
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
