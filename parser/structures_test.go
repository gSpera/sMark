package parser

import (
	"eNote/token"
	"errors"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestStructure(t *testing.T) {
	tests := []struct {
		name   string
		input  []token.Token
		output []token.Token
	}{
		{
			"Bold I",
			[]token.Token{
				token.BoldToken{},
				token.TextToken{Text: "Test"},
				token.BoldToken{},
			},
			[]token.Token{
				token.TextToken{Text: "Test", Bold: true},
			},
		}, //End Bold I
		{
			"Bold II",
			[]token.Token{
				token.BoldToken{},
				token.TextToken{Text: "Test "},
				token.TextToken{Text: "Another Text"},
				token.BoldToken{},
			},
			[]token.Token{
				token.TextToken{Text: "Test Another Text", Bold: true},
			},
		}, //End Bold II,
		{
			"Bold Not Ending",
			[]token.Token{
				token.BoldToken{},
				token.TextToken{Text: "Test"},
			},
			[]token.Token{
				token.BoldToken{},
				token.TextToken{Text: "Test"}},
		}, //End Bold Not Ending
		{
			"Italic",
			[]token.Token{
				token.ItalicToken{},
				token.TextToken{Text: "Test"},
				token.ItalicToken{},
			},
			[]token.Token{
				token.TextToken{Text: "Test", Italic: true},
			},
		},
		{
			"Italic Not Ending",
			[]token.Token{
				token.ItalicToken{},
				token.TextToken{Text: "Test"},
			},
			[]token.Token{
				token.ItalicToken{},
				token.TextToken{Text: "Test"}},
		}, //End Italic Not Ending
		{
			"Strike-Throught",
			[]token.Token{
				token.LessToken{},
				token.TextToken{Text: "Test"},
				token.LessToken{},
			},
			[]token.Token{
				token.TextToken{Text: "Test", Strike: true},
			},
		},
		{
			"Strike-Throught Not Ending",
			[]token.Token{
				token.LessToken{},
				token.TextToken{Text: "Test"},
			},
			[]token.Token{
				token.LessToken{},
				token.TextToken{Text: "Test"}},
		}, //End Italic Not Ending
		{
			"Token after end",
			[]token.Token{
				token.ItalicToken{},
				token.TextToken{Text: "Test"},
				token.ItalicToken{},
				token.TextToken{Text: "AnotherText"},
			},
			[]token.Token{
				token.TextToken{Text: "Test", Italic: true},
				token.TextToken{Text: "AnotherText"},
			},
		},
		{
			"CheckBox",
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: " "},
				token.SBracketCloseToken{},
			},
			[]token.Token{
				token.CheckBoxToken{Char: ' '},
			},
		},
		{
			"CheckBox_NotEnoughTokens",
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: " "},
			},
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: " "},
			},
		},
		{
			"CheckBox_NotClose",
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: " "},
				token.BoldToken{},
			},
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: " "},
				token.BoldToken{},
			},
		},
		{
			"CheckBox_StringNoCorrectSize",
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: "  "},
				token.SBracketCloseToken{},
			},
			[]token.Token{
				token.SBracketOpenToken{},
				token.TextToken{Text: "  "},
				token.SBracketCloseToken{}},
		},
		{
			"CheckBox_StringNoCorrectSize",
			[]token.Token{
				token.SBracketOpenToken{},
				token.BoldToken{},
				token.SBracketCloseToken{},
			},
			[]token.Token{
				token.SBracketOpenToken{},
				token.BoldToken{},
				token.SBracketCloseToken{}},
		},
		{
			"Two Times",
			[]token.Token{
				token.ItalicToken{},
				token.TextToken{Text: "Test"},
				token.ItalicToken{},
				token.ItalicToken{},
				token.TextToken{Text: "AnotherText"},
				token.ItalicToken{},
			},
			[]token.Token{
				token.TextToken{Text: "Test", Italic: true},
				token.TextToken{Text: "AnotherText", Italic: true},
			},
		},
		{
			"Adiacent",
			[]token.Token{
				token.BoldToken{},
				token.BoldToken{},
			},
			[]token.Token{
				token.BoldToken{},
				token.BoldToken{},
			},
		},
		{
			"Token inside",
			[]token.Token{
				token.BoldToken{},
				token.CheckBoxToken{Char: 'X'},
				token.BoldToken{},
			},
			[]token.Token{
				token.BoldToken{},
				token.CheckBoxToken{Char: 'X'},
				token.BoldToken{},
			},
		},
	}

	spew.Config = spew.ConfigState{
		DisableMethods: true,
		Indent:         "\t",
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			structure := TokenToStructure(test.input)

			if ok, err := checkSlice(structure, test.output); !ok {
				fmt.Println("Not expected output")
				fmt.Println(err)

				fmt.Println("Expected:", spew.Sdump(test.output))
				fmt.Print("Got:", spew.Sdump(structure))
				t.Fail()
			}
		})
	}
}

func TestSearchCheckBoxPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected panic")
		}
	}()

	searchCheckbox([]token.Token{
		token.BoldToken{},
		token.BoldToken{},
		token.BoldToken{},
	}, 0)
}

//checkSlice checks if two tokens slice are equal,
//it checks for len and elements
func checkSlice(arr1, arr2 []token.Token) (bool, error) {
	if len(arr1) != len(arr2) {
		return false, errors.New("Len different")
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false, fmt.Errorf("Differ at pos: %d", i)
		}
	}

	return true, nil
}
