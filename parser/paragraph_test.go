package parser

import (
	"eNote/token"
	"testing"
)

func TestNotEmptyLines(t *testing.T) {
	tm := []struct {
		name   string
		input  []token.LineContainer
		output uint
	}{
		{
			"Empty",
			[]token.LineContainer{},
			0,
		},
		{
			"One",
			[]token.LineContainer{
				token.LineContainer{
					Tokens: []token.Token{
						token.BoldToken{},
					},
				},
			},
			1,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := notEmptyLines(tt.input)
			if tt.output != out {
				t.Errorf("Expected: %d; got: %d", tt.output, out)
			}
		})
	}
}

func TestEqualLine(t *testing.T) {
	tm := []struct {
		name      string
		paragraph token.TextParagraph
		lines     []token.LineToken
		index     int
		output    token.ParagraphToken
		outputOk  bool
	}{
		{
			"Empty",
			token.TextParagraph{},
			[]token.LineToken{},
			0,
			nil,
			false,
		},
		{
			"Index 0",
			token.TextParagraph{},
			[]token.LineToken{
				token.EqualLine{},
				token.LineContainer{},
			},
			0,
			nil,
			false,
		},
		{
			"Index not EqualLine",
			token.TextParagraph{},
			[]token.LineToken{
				token.LineContainer{},
				token.LineContainer{},
			},
			0,
			nil,
			false,
		},
		{
			"Index 0",
			token.TextParagraph{},
			[]token.LineToken{
				token.LineContainer{},
				token.EqualLine{},
			},
			1,
			nil,
			false,
		},
		{
			"Not Line Container",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{
						token.BoldToken{},
					}},
				},
			},
			[]token.LineToken{
				token.EqualLine{},
				token.EqualLine{},
			},
			1,
			nil,
			false,
		},
		{
			"Not Line Container",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{
						token.BoldToken{},
					}}},
			},
			[]token.LineToken{
				token.LineContainer{},
				token.EqualLine{},
			},
			1,
			nil,
			false,
		},
		{
			"Not Line Container",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{
						token.BoldToken{},
					}}},
			},
			[]token.LineToken{
				token.LineContainer{
					Tokens: []token.Token{
						token.BoldToken{},
					},
				},
				token.EqualLine{Indentation: 1},
			},
			1,
			nil,
			false,
		},
		{
			"Title",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{
						token.BoldToken{},
					}}},
			},
			[]token.LineToken{
				token.LineContainer{
					Indentation: 1,
					Tokens: []token.Token{
						token.BoldToken{},
					},
				},
				token.EqualLine{Indentation: 1},
			},
			1,
			token.TitleParagraph{
				Text:        "*",
				Indentation: 1,
			},
			true,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out, ok := equalLine(tt.paragraph, tt.lines, tt.index)
			if tt.outputOk != ok {
				t.Errorf("Not OK: Expected %v; got: %v", tt.outputOk, ok)
			}

			if tt.output != out {
				t.Errorf("Expected: %v; got: %v", tt.output, out)
			}
		})
	}
}
