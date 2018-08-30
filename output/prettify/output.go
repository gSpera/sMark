package prettify

import (
	"fmt"
	"log"
	"strings"

	"github.com/gSpera/sMark/token"
	sMark "github.com/gSpera/sMark/utils"
)

const divisorLen = 10

//Output generate sMark valid code as output, it can be used to prettify the sMark source code
func Output(paragraphs []token.ParagraphToken, options sMark.Options) ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, makeHeader(options)...)

	for _, p := range paragraphs {
		data = append(data, []byte(Single(p))...)
	}
	return data, nil
}

//Single generate the sMark code for a single paragraph
func Single(paragraph token.ParagraphToken) string {
	switch pp := paragraph.(type) {
	case token.TitleParagraph:
		return makeTitle(pp, "=")
	case token.SubtitleParagraph:
		return makeTitle(token.TitleParagraph(pp), "-")
	case token.TextParagraph:
		return textParagraphToString(pp)
	case token.DivisorParagraph:
		return fmt.Sprintf("%s\n\n", strings.Repeat("-", divisorLen))
	case token.ListParagraph:
		return makeList(pp)
	case token.CodeParagraph:
		return makeCode(pp)
	default:
		log.Fatalf("\t\t- ERROR: Not Implemented: %T{%v}\n", pp, pp)
		return ""
	}
}

func textParagraphToString(p token.TextParagraph) string {
	var str string
	for _, line := range p.Lines {
		lineString := strings.Repeat("\t", line.Indentation)
		lineString += lineContainerToString(line)
		if len(lineString) == 0 { //Do not add line if it is empty
			continue
		}
		str += lineString + "\n"
	}
	if len(str) != 0 {
		str += "\n"
	}
	return str
}

func makeTitle(title token.TitleParagraph, char string) string {
	var str string
	str = fmt.Sprintf("%s%s\n%s%s\n\n",
		indentation(title.Indentation), title.Text,
		indentation(title.Indentation), strings.Repeat(char, len(title.Text)),
	)
	return str
}

func makeHeader(options sMark.Options) string {
	var maxLen int
	var content string

	for k, v := range options.String {
		line := fmt.Sprintf("%s = %s\n", k, v)
		maxLen = max(maxLen, len(line))
		content += line
	}
	for k, v := range options.Bool {
		line := fmt.Sprintf("%s = %t\n", k, v)
		maxLen = max(maxLen, len(line))
		content += line
	}
	for k, v := range options.Generic {
		line := fmt.Sprintf("%s = %v\n", k, v)
		maxLen = max(maxLen, len(line))
		content += line
	}

	delim := strings.Repeat("+", maxLen)

	return fmt.Sprintf("%s\n%s%s\n\n", delim, content, delim)
}

func indentation(n int) string {
	return strings.Repeat("\t", int(n))
}

func makeList(list token.ListParagraph) string {
	var content string
	for _, item := range list.Items {
		content += fmt.Sprintf("%s %s\n", strings.Repeat("-", item.Indentation), item.Text.String())
	}

	content += "\n"
	return content
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func lineContainerToString(container token.LineContainer) string {
	var str string
	if container.Quote {
		str = "|"
	}

	for _, t := range container.Tokens {
		switch tt := t.(type) {
		case token.TextToken:
			str += tt.StringEscape()
		case token.SimpleToken:
			str += string(tt.Char())
		case token.CheckBoxToken:
			str += fmt.Sprintf("[%c]", tt.Char)
		default:
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", tt, tt))
		}
	}
	return str
}

func makeCode(p token.CodeParagraph) string {
	var str string
	str = fmt.Sprintf("[%s]\n", p.Lang)
	str += textParagraphToString(p.Text)
	str += "[end]\n"
	return str
}
