package token

import (
	"fmt"
	"testing"
)

type badToken struct{ Token }

func (t badToken) Type() Type { return -1 }

func TestLine_String(t *testing.T) {
	tt := []struct {
		Name   string
		Input  LineContainer
		Result string
	}{
		{"Simple", LineContainer{Indentation: 0, Tokens: []Token{TextToken{Text: "Test"}}}, "Test"},
		{"Indentation", LineContainer{Indentation: 1, Tokens: []Token{TextToken{Text: "Test"}}}, "\tTest"},
		{"SimpleToken", LineContainer{Indentation: 0, Tokens: []Token{BoldToken{}}}, "*"},
		{"CheckBox", LineContainer{Indentation: 0, Tokens: []Token{CheckBoxToken{Char: ' '}}}, "[ ]"},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			result := test.Input.String()
			if test.Result != result {
				t.Fatalf("Strings are not equal: got %s, expected: %s", result, test.Result)
			}
		})
	}

	t.Run("UnknownToken_Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Panic is not occured")
			}
		}()

		result := LineContainer{Indentation: 0, Tokens: []Token{badToken{}}}.String()
		fmt.Println("Result:", result)
	})
}
