package parser

import (
	"eNote/token"
	"eNote/utils"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestTokenToLine(t *testing.T) {
	tm := []struct {
		name   string
		input  []token.Token
		output []token.LineToken
	}{
		{
			"No Input",
			[]token.Token{},
			[]token.LineToken{},
		},
		{
			"Simple Text",
			[]token.Token{
				token.TextToken{Text: "Text"},
				token.NewLineToken{},
			},
			[]token.LineToken{
				token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Text"}}},
			},
		},
		{
			"No NewLine",
			[]token.Token{
				token.TextToken{Text: "Text"},
			},
			[]token.LineToken{
				token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Text"}}},
			},
		},
		{
			"Tab Indentation",
			[]token.Token{
				token.TabToken{},
				token.TabToken{},
				token.TextToken{Text: "Text"},
			},
			[]token.LineToken{
				token.LineContainer{
					Indentation: 2,
					Tokens:      []token.Token{token.TextToken{Text: "Text"}},
				},
			},
		},
		{
			"Tab No Indentation",
			[]token.Token{
				token.TabToken{},
				token.TextToken{Text: "Text"},
				token.TabToken{},
			},
			[]token.LineToken{
				token.LineContainer{
					Indentation: 1,
					Tokens: []token.Token{
						token.TextToken{Text: "Text"},
						token.TabToken{},
					},
				},
			},
		},
	}
	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			result := TokenToLine(tt.input)
			ok, err := checkLineSlice(tt.output, result)
			if !ok {
				t.Error(err)
			}
		})
	}
}

//checkLineSlice checks if two line tokens slice are equal,
//it checks for len and elements
func checkLineSlice(arr1, arr2 []token.LineToken) (bool, error) {
	if len(arr1) != len(arr2) {
		return false, errors.New("Len different")
	}

	for i := range arr1 {
		c1, ok1 := arr1[i].(token.LineContainer)
		c2, ok2 := arr2[i].(token.LineContainer)
		if ok1 != ok2 {
			return false, fmt.Errorf("Differ at pos: %d, only one LineContainer found", i)
		} else if ok1 {
			return checkLineContainer(c1, c2)
		}

		if arr1[i] != arr2[i] {
			return false, fmt.Errorf("Differ at pos: %d", i)
		}
	}

	return true, nil
}

func checkLineContainer(l1, l2 token.LineContainer) (bool, error) {
	if l1.Indentation != l2.Indentation {
		return false, fmt.Errorf("LineContainer: Indentation is different")
	}
	if l1.Quote != l2.Quote {
		return false, fmt.Errorf("LineContainer: Quote is different")
	}

	return checkSlice(l1.Tokens, l2.Tokens)
}

func TestCheckIndentation(t *testing.T) {
	tm := []struct {
		name   string
		input  token.TextParagraph
		output int
	}{
		{
			"Simple",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Test"}}},
				},
			},
			0,
		},
		{
			"Empty",
			token.TextParagraph{
				Indentation: 1,
				Lines:       []token.LineContainer{},
			},
			0,
		},
		{
			"Not Constant",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{
						Indentation: 3,
					},
					token.LineContainer{
						Indentation: 2,
					},
				},
			},
			0,
		},
		{
			"Right",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{
						Indentation: 3,
					},
					token.LineContainer{
						Indentation: 3,
					},
				},
			},
			3,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := &tt.input
			checkIndentation(out)
			if tt.output != out.Indentation {
				spew.Dump(tt.output)
				spew.Dump(*out)
				t.Errorf("Expected: %v; got: %v", tt.output, out.Indentation)
			}
		})
	}
}

func checkTextParagraph(p1, p2 token.TextParagraph) (bool, error) {
	if p1.Indentation != p2.Indentation {
		return false, fmt.Errorf("Indentation is different")
	}
	if len(p1.Lines) != len(p2.Lines) {
		return false, fmt.Errorf("Len is different")
	}

	for i := range p1.Lines {
		if ok, err := checkLineContainer(p1.Lines[i], p2.Lines[i]); !ok {
			return false, err
		}
	}
	return true, nil
}

func TestIsOnlyWhiteSpace(t *testing.T) {
	tm := []struct {
		name   string
		input  string
		output bool
	}{
		{
			"Empty",
			"",
			true,
		},
		{
			"WhiteSpace",
			" ",
			true,
		},
		{
			"Tab",
			"\t",
			true,
		},
		{
			"Text",
			"test",
			false,
		},
		{
			"Text and whitespace",
			" test ",
			false,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := isOnlyWhiteSpace(tt.input)
			if tt.output != out {
				t.Errorf("Expected %v; got: %v", tt.output, out)
			}
		})
	}
}

func TestIsQuoteLine(t *testing.T) {
	tm := []struct {
		name   string
		input  token.LineContainer
		output bool
	}{
		{
			"Empty",
			token.LineContainer{},
			false,
		},
		{
			"PipeToken",
			token.LineContainer{
				Tokens: []token.Token{
					token.PipeToken{},
				},
			},
			true,
		},
		{
			"OtherToken",
			token.LineContainer{
				Tokens: []token.Token{
					token.BoldToken{},
				},
			},
			false,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := isQuoteLine(tt.input)
			if tt.output != out {
				t.Errorf("Expected: %v; got: %v", tt.output, out)
			}
		})
	}
}

func TestIsType(t *testing.T) {
	tm := []struct {
		name    string
		inType  token.Type
		line    token.LineContainer
		options isTypeOptions
		output  bool
	}{
		{
			"Empty",
			token.TypeBold,
			token.LineContainer{},
			isTypeOptions{},
			false,
		},
		{
			"Bold",
			token.TypeBold,
			token.LineContainer{
				Tokens: []token.Token{
					token.BoldToken{},
					token.BoldToken{},
				},
			},
			isTypeOptions{},
			true,
		},
		{
			"Ignore Tag",
			token.TypeBold,
			token.LineContainer{
				Tokens: []token.Token{
					token.TabToken{},
					token.BoldToken{},
					token.BoldToken{},
				},
			},
			isTypeOptions{
				ignoreTabs: true,
			},
			true,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := isType(tt.inType, tt.line, []isTypeOptions{tt.options}...)
			if tt.output != out {
				t.Errorf("Expected: %v; got: %v", tt.output, out)
			}
		})
	}
}

func TestIsCodeHeader(t *testing.T) {
	tm := []struct {
		name   string
		input  token.LineContainer
		output bool
	}{
		{
			"Empty",
			token.LineContainer{},
			false,
		},
		{
			"Not Opening",
			token.LineContainer{
				Tokens: []token.Token{
					token.BoldToken{},
					token.BoldToken{},
					token.BoldToken{},
				},
			},
			false,
		},
		{
			"Not Close",
			token.LineContainer{
				Tokens: []token.Token{
					token.SBracketOpenToken{},
					token.BoldToken{},
					token.BoldToken{},
				},
			},
			false,
		},
		{
			"Not Text",
			token.LineContainer{
				Tokens: []token.Token{
					token.SBracketOpenToken{},
					token.BoldToken{},
					token.SBracketCloseToken{},
				},
			},
			false,
		},
		{
			"Code Header",
			token.LineContainer{
				Tokens: []token.Token{
					token.SBracketOpenToken{},
					token.TextToken{Text: "Test"},
					token.SBracketCloseToken{},
				},
			},
			true,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := isCodeHeader(tt.input)
			if tt.output != out {
				t.Errorf("Expected: %v; got: %v", tt.output, out)
			}
		})
	}
}

func TestParseHeaderLine(t *testing.T) {
	tm := []struct {
		name   string
		input  token.TextParagraph
		output eNote.Options
	}{
		{
			"Empty",
			token.TextParagraph{},
			eNote.Options{},
		},
		{
			"Value",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Key = Value"}}},
				},
			},
			eNote.Options{
				String: map[string]string{
					"Key": "Value",
				},
			},
		},
		{
			"Comment",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Key = Value; This is a comment"}}},
				},
			},
			eNote.Options{
				String: map[string]string{
					"Key": "Value",
				},
			},
		},
		{
			"White Space",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Ke y = Value; This is a comment"}}},
				},
			},
			eNote.Options{
				String: map[string]string{
					"Ke y": "Value",
				},
			},
		},
		{
			"White Space At End",
			token.TextParagraph{
				Lines: []token.LineContainer{
					token.LineContainer{Tokens: []token.Token{token.TextToken{Text: "Key = Value "}}},
				},
			},
			eNote.Options{
				String: map[string]string{
					"Key": "Value",
				},
			},
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := parseHeaderLines(tt.input)
			if ok, err := checkOptions(tt.output, out); !ok {
				t.Error(err)
			}
		})
	}
}

func checkOptions(o1, o2 eNote.Options) (bool, error) {
	if len(o1.Bool) != len(o2.Bool) {
		return false, fmt.Errorf("Bool: Len is different")
	}
	if len(o1.String) != len(o2.String) {
		return false, fmt.Errorf("String: Len is different")
	}
	if len(o1.Generic) != len(o2.Generic) {
		return false, fmt.Errorf("Generic: Len is different")
	}

	for key := range o2.Bool {
		if o1.Bool[key] != o2.Bool[key] {
			return false, fmt.Errorf("Bool: Different at key: [%s ]: [%v] != [%v]", key, o1.Bool[key], o2.Bool[key])
		}
	}
	for key := range o2.String {
		if o1.String[key] != o2.String[key] {
			return false, fmt.Errorf("String: Different at key: [%s ]: [%v] != [%v]", key, o1.String[key], o2.String[key])
		}
	}
	for key := range o2.Generic {
		if o1.Generic[key] != o2.Generic[key] {
			return false, fmt.Errorf("Generic: Different at key: [%s ]: [%v] != [%v]", key, o1.Generic[key], o2.Generic[key])
		}
	}

	return true, nil
}

func TestParseLine(t *testing.T) {
	tm := []struct {
		name      string
		input     token.LineContainer
		output    token.LineToken
		deepCheck bool
	}{
		{
			"Empty",
			token.LineContainer{},
			token.LineContainer{},
			true,
		},
		{
			"HeaderLine",
			token.LineContainer{
				Tokens: []token.Token{
					token.HeaderToken{},
					token.HeaderToken{},
				},
			},
			token.HeaderLine{},
			false,
		},
		{
			"EqualLine",
			token.LineContainer{
				Tokens: []token.Token{
					token.EqualToken{},
					token.EqualToken{},
				},
			},
			token.EqualLine{
				Length: 2,
			},
			true,
		},
		{
			"EqualLine Tab",
			token.LineContainer{
				Tokens: []token.Token{
					token.TabToken{},
					token.EqualToken{},
					token.EqualToken{},
				},
			},
			token.EqualLine{
				Length: 3,
			},
			true,
		},
		{
			"LessLine",
			token.LineContainer{
				Tokens: []token.Token{
					token.LessToken{},
					token.LessToken{},
				},
			},
			token.LessLine{},
			false,
		},
		{
			"ListLine",
			token.LineContainer{
				Tokens: []token.Token{
					token.LessToken{},
					token.TextToken{
						Text: "Test",
					},
				},
			},
			token.ListLine{},
			false,
		},
		{
			"CodeLine",
			token.LineContainer{
				Tokens: []token.Token{
					token.SBracketOpenToken{},
					token.TextToken{Text: "Test"},
					token.SBracketCloseToken{},
				},
			},
			token.CodeLine{Lang: "Test"},
			true,
		},
		{
			"QuoteLine",
			token.LineContainer{
				Tokens: []token.Token{
					token.PipeToken{},
					token.TextToken{Text: "Test"},
				},
			},
			token.LineContainer{Tokens: []token.Token{
				token.TextToken{Text: "Test"},
			}, Quote: true},
			true,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := parseLine(tt.input)
			lc1, ok1 := tt.output.(token.LineContainer)
			lc2, ok2 := out.(token.LineContainer)
			if ok1 && ok2 {
				ok, err := checkLineContainer(lc1, lc2)
				if !ok {
					t.Error(err)
				}
				return
			}
			if !tt.deepCheck {
				t1 := reflect.TypeOf(tt.output)
				t2 := reflect.TypeOf(out)
				if t1 != t2 {
					t.Error("Types are different")
				}
			}
			if tt.deepCheck && tt.output != out {
				t.Errorf("Expected %T{%+v}; got: %T{%+v}", tt.output, tt.output, out, out)
			}
		})
	}
}

func TestIsListLine(t *testing.T) {
	tm := []struct {
		name   string
		line   token.LineContainer
		output bool
	}{
		{
			"Empty",
			token.LineContainer{},
			false,
		},
		{
			"False",
			token.LineContainer{
				Tokens: []token.Token{
					token.BoldToken{},
					token.BoldToken{},
				},
			},
			false,
		},
		{
			"True",
			token.LineContainer{
				Tokens: []token.Token{
					token.LessToken{},
					token.BoldToken{},
				},
			},
			true,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := isListLine(tt.line)
			if tt.output != out {
				t.Errorf("Expected %t; got: %t", tt.output, out)
			}
		})
	}
}
