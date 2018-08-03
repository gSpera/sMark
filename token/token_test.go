package token

import "testing"

func TestTextToken_String(t *testing.T) {
	tt := []struct {
		Name   string
		Input  TextToken
		Result string
	}{
		{"Simple Text", TextToken{Text: "Test"}, "Test"},
		{"Bold Text", TextToken{Text: "Test", Bold: true}, "*Test*"},
		{"Italic Text", TextToken{Text: "Test", Italic: true}, "/Test/"},
		{"Strike-Throught Text", TextToken{Text: "Test", Strike: true}, "-Test-"},
		{"Indentation Text", TextToken{Text: "Test", Indentation: 3}, "\t\t\tTest"},
		{"Bold and Italic Text", TextToken{Text: "Test", Bold: true, Italic: true}, "*/Test*/"},
		{"Link", TextToken{Text: "Test", Link: "https://example.com"}, "\"Test\"@\"https://example.com\""},
	}

	for _, test := range tt {
		t.Run(test.Name+"_String", func(t *testing.T) {
			result := test.Input.String()
			if result != test.Result {
				t.Fatalf("Strings are not equal: got %s, expected: %s", result, test.Result)
			}
		})
		t.Run(test.Name+"_EscapeString", func(t *testing.T) {
			result := test.Input.StringEscape()
			if result != test.Result {
				t.Fatalf("Strings are not equal: got %s, expected: %s", result, test.Result)
			}
		})
	}

}
