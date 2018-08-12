package parser

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/gSpera/sMark/token"

	"github.com/davecgh/go-spew/spew"
)

type badReader struct{}

func (b badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("You are reading from badReader")
}

func TestTokenizer(t *testing.T) {
	tt := []struct {
		Name   string
		Data   string
		Result []token.Token
		Err    error
	}{
		{"Bold Token", "*Test*\n", []token.Token{
			token.BoldToken{}, token.TextToken{Text: "Test"}, token.BoldToken{}, token.NewLineToken{},
		}, nil},
		{"Italic Token", "/Test/\n", []token.Token{
			token.ItalicToken{}, token.TextToken{Text: "Test"}, token.ItalicToken{}, token.NewLineToken{},
		}, nil},
		{"Strike-Throught Token", "-Test-\n", []token.Token{
			token.LessToken{}, token.TextToken{Text: "Test"}, token.LessToken{}, token.NewLineToken{},
		}, nil},
		{"NewLine Token", "Test\nTest\n", []token.Token{
			token.TextToken{Text: "Test"}, token.NewLineToken{}, token.TextToken{Text: "Test"}, token.NewLineToken{},
		}, nil},
		{"Empty Text", "", []token.Token{}, nil},
		{"Escape", "\\*Test\\*\n", []token.Token{
			token.TextToken{Text: "*Test*"}, token.NewLineToken{},
		}, nil},
		{"Auto NewLine", "*", []token.Token{
			token.BoldToken{}, token.NewLineToken{},
		}, nil},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			in := bytes.NewBufferString(test.Data)
			toks, err := Tokenizer(in)
			if (test.Err != nil) != (err != nil) {
				log.Println("Expected:", test.Err)
				log.Println("Got:", err)
				t.Error("Error is not the same")
			}
			if !checkResult(toks, test.Result) {
				fmt.Println("Got:")
				spew.Dump(toks)
				fmt.Println("Expected:")
				spew.Dump(test.Result)
				t.Errorf("Result is not the same: got: %v; expected: %v", toks, test.Result)
			}
		})
	}

	t.Run("Panic", func(t *testing.T) {
		_, err := Tokenizer(badReader{})
		if err == nil {
			t.Fatalf("Could not tokenize: %v", err)
		}
	})
}

func TestIn(t *testing.T) {
	tt := []struct {
		Name   string
		Token  token.Type
		List   []token.Type
		Result bool
	}{
		{"True", token.TypeBold, []token.Type{token.TypeItalic, token.TypeBold}, true},
		{"False", token.TypeBold, []token.Type{token.TypeItalic, token.TypeLess}, false},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			result := in(test.Token, test.List)
			if test.Result != result {
				t.Fatalf("In is not correct: got %v, expcted: %v", result, test.Result)
			}
		})
	}
}

func TestCheckTokenList(t *testing.T) {
	tt := []struct {
		Name   string
		Input  []token.Token
		Result []token.Token
	}{
		{"Empty", []token.Token{}, []token.Token{}},
		{"RemoveFirst", []token.Token{token.TextToken{}}, []token.Token{}},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			result := checkTokenList(test.Input)
			if !checkResult(result, test.Result) {
				t.Fatalf("In is not correct: got %v, expcted: %v", result, test.Result)
			}
		})
	}
}

func checkResult(t, tt []token.Token) bool {
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
