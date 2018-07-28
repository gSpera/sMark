package token

import (
	"testing"
)

func TestType(t *testing.T) {
	t.Parallel()
	tt := []struct {
		Name   string
		Token  Token
		Result Type
	}{
		{"BoldToken", BoldToken{}, TypeBold},
		{"ItalicToken", ItalicToken{}, TypeItalic},
		{"BoldToken", BoldToken{}, TypeBold},
		{"ItalicToken", ItalicToken{}, TypeItalic},
		{"NewLineToken", NewLineToken{}, TypeNewLine},
		{"TabToken", TabToken{}, TypeTab},
		{"HeaderToken", HeaderToken{}, TypeHeader},
		{"LessToken", LessToken{}, TypeLess},
		{"EqualToken", EqualToken{}, TypeEqual},
		{"TextToken", TextToken{}, TypeText},
		{"HeaderParagraph", HeaderParagraph{}, TypeParagraphHeader},
		{"TextParagraph", TextParagraph{}, TypeParagraphText},
		{"DivisorParagraph", DivisorParagraph{}, TypeParagraphDivisor},
		{"TitleParagraph", TitleParagraph{}, TypeParagraphTitle},
		{"SubtitleParagraph", SubtitleParagraph{}, TypeParagraphSubtitle},
		{"ListParagraph", ListParagraph{}, TypeParagraphList},
		{"SBracketOpenToken", SBracketOpenToken{}, TypeSBracketOpen},
		{"SBracketCloseToken", SBracketCloseToken{}, TypeSBracketClose},
		{"CheckBoxToken", CheckBoxToken{}, TypeCheckBox},
		{"EscapeToken", EscapeToken{}, TypeEscape},
		{"PipeToken", PipeToken{}, TypePipe},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			result := test.Token.Type()
			if test.Result != result {
				t.Fatalf("Type is not the same: got %v; expected %v", result, test.Result)
			}
		})
	}
}
