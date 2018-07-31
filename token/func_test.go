package token

import "testing"

func TestSimpleTokenChar(t *testing.T) {
	t.Parallel()
	tt := []struct {
		Name   string
		Token  SimpleToken
		Result rune
	}{
		{"BoldToken", BoldToken{}, '*'},
		{"ItalicToken", ItalicToken{}, '/'},
		{"HeaderToken", HeaderToken{}, '+'},
		{"NewLineToken", NewLineToken{}, '\n'},
		{"TabToken", TabToken{}, '\t'},
		{"LessToken", LessToken{}, '-'},
		{"EqualToken", EqualToken{}, '='},
		{"SBracketOpenToken", SBracketOpenToken{}, '['},
		{"SBracketCloseToken", SBracketCloseToken{}, ']'},
		{"PipeToken", PipeToken{}, '|'},
		{"QuoteToken", QuoteToken{}, '"'},
		{"AtToken", AtToken{}, '@'},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			result := test.Token.Char()
			if test.Result != result {
				t.Fatalf("Chars are different: got %v; expected: %v", result, test.Result)
			}
		})
	}
}
