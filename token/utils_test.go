package token

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestLineContainer(t *testing.T) {
	tt := []struct {
		Name   string
		String string
		Result LineContainer
	}{
		{
			"No Tab", "Test",
			LineContainer{
				Tokens:      []Token{TextToken{Text: "Test"}},
				Indentation: 0,
			},
		},
		{
			"Tab", "\tTest",
			LineContainer{
				Tokens:      []Token{TextToken{Text: "Test"}},
				Indentation: 1,
			},
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			line := LineContainerFromString(test.String)
			if test.Result.Indentation != line.Indentation {
				t.Fatal("Indentation is not the same")
			}
			if !checkResult(line.Tokens, test.Result.Tokens) {
				spew.Dump(line.Tokens, test.Result.Tokens)
				t.Fatal("Result is not the same")
			}
		})
	}
}

func checkResult(t, tt []Token) bool {
	if len(t) != len(tt) {
		return false
	}
	for i := range tt {
		if t[i] != tt[i] {
			return false
		}
	}

	return true
}
