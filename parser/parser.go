package parser

import (
	"eNote/token"
	"io"

	"github.com/davecgh/go-spew/spew"
)

//ParseReader parse a source file and returns the paragraphs
func ParseReader(fl io.Reader) ([]token.ParagraphToken, error) {
	tokens, err := Tokenizer(fl)
	spew.Dump(tokens)
	if err != nil {
		return nil, err
	}

	lines := TokenToLine(tokens)
	paragraphs := TokenToParagraph(lines)

	return paragraphs, nil
}
