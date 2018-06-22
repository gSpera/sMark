package prettify

import (
	"eNote/token"
	eNote "eNote/utils"
	"fmt"
	"log"
	"strings"
)

const divisorLen = 10

//Output generate eNote valid code as output, it can be used to prettify the eNote source code
func Output(paragraphs []token.ParagraphToken, options eNote.Options) ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, makeHeader(options)...)

	for _, p := range paragraphs {
		switch pp := p.(type) {
		case token.TitleParagraph:
			data = append(data, makeTitle(pp)...)
		case token.SubtitleParagraph:
			data = append(data, makeSubtitle(pp)...)
		case token.TextParagraph:
			data = append(data, textParagraphToString(pp)...)
		case token.DivisorParagraph:
			data = append(data, fmt.Sprintf("%s\n\n", strings.Repeat("-", divisorLen))...)
		case token.ListParagraph:
			data = append(data, makeList(pp)...)
		default:
			log.Printf("\t\t- ERROR: Not Implemented: %T{%v}\n", pp, pp)
		}
	}
	return data, nil
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

func makeTitle(title token.TitleParagraph) string {
	var str string
	str = fmt.Sprintf("%s%s\n%s%s\n\n",
		indentation(title.Indentation), title.Text,
		indentation(title.Indentation), strings.Repeat("=", len(title.Text)-1),
	)
	return str
}

func makeSubtitle(subtitle token.SubtitleParagraph) string {
	var str string
	str = fmt.Sprintf("%s%s\n%s%s\n\n",
		indentation(subtitle.Indentation), subtitle.Text,
		indentation(subtitle.Indentation), strings.Repeat("-", len(subtitle.Text)-1),
	)
	return str
}

func makeHeader(options eNote.Options) string {
	delim := strings.Repeat("+", 10)
	var content string

	content += fmt.Sprintf("Title=%v\n", strings.TrimSpace(*options.Title))
	content += fmt.Sprintf("TabWidth=%v\n", *options.TabWidth)
	content += fmt.Sprintf("InputFile=%v\n", *options.InputFile)
	content += fmt.Sprintf("NewLine=%v\n", *options.NewLine)
	content += fmt.Sprintf("OutputFile=%v\n", *options.OutputFile)
	content += fmt.Sprintf("OnlyBody=%v\n", *options.OnlyBody)
	content += fmt.Sprintf("EnableFont=%v\n", *options.EnableFont)
	content += fmt.Sprintf("InlineCSS=%v\n", *options.InlineCSS)
	content += fmt.Sprintf("CustomCSS=%v\n", *options.CustomCSS)

	return fmt.Sprintf("%s\n%s%s\n\n", delim, content, delim)
}

func indentation(n int) string {
	return strings.Repeat("\t", int(n))
}

func makeList(list token.ListParagraph) string {
	var content string
	for _, item := range list.Items {
		content += fmt.Sprintf("- %s\n", item.Text.String())
	}

	content += "\n"
	return content
}

func lineContainerToString(container token.LineContainer) string {
	var str string
	for _, t := range container.Tokens {
		switch tt := t.(type) {
		case token.TextToken:
			str += textTokenToString(tt)
		default:
			if st, ok := tt.(token.SimpleToken); ok {
				str += string(st.Char())
				continue
			}
			panic(fmt.Sprintf("LineContainer contains unknown token %T{%v}", tt, tt))
		}
	}
	return str
}

func textTokenToString(t token.TextToken) string {
	return t.String()
}
