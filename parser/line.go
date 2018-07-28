package parser

import (
	"eNote/token"
	eNote "eNote/utils"
	"log"
	"strings"
)

//TokenToLine divide a slice of tokens in lines
func TokenToLine(tokens []token.Token) []token.LineToken {
	lines := []token.LineToken{}
	currentLine := token.LineContainer{}

	indent := true

	for _, t := range tokens {
		switch t.(type) {
		case token.TabToken:
			if indent {
				currentLine.Indentation++
			} else {
				currentLine.Tokens = append(currentLine.Tokens, token.TabToken{})
				indent = false
			}
		case token.NewLineToken:
			line := parseLine(currentLine)
			lines = append(lines, line)
			currentLine = token.LineContainer{}
			indent = true
		default:
			currentLine.Tokens = append(currentLine.Tokens, t)
			indent = false
		}

	}

	return lines
}

const isTypeThreshold = 2

type isTypeOptions struct {
	threshold  int
	ignoreTabs bool
}

//isType reutnrs wheter or not the passed line contains only the specified type of tokens
//if the line contains less than the threshold tokens it will return always false
//the default value for threshold is isTypeThreshold contant
func isType(typ token.Type, line token.LineContainer, _options ...isTypeOptions) bool {
	//Calculate the threshold
	options := isTypeOptions{
		threshold:  isTypeThreshold,
		ignoreTabs: false,
	}

	if len(_options) != 0 {
		options = _options[0]
	}
	if options.threshold == 0 {
		options.threshold = 2
	}
	//Check if there are too few tokens
	if len(line.Tokens) < options.threshold {
		return false
	}

	//Checks the line
	for _, t := range line.Tokens {
		if options.ignoreTabs && t.Type() == token.TypeTab {
			continue
		}
		if typ != t.Type() {
			return false
		}
	}

	return true
}

func isListLine(line token.LineContainer) bool {
	if len(line.Tokens) < 2 {
		return false
	}

	if _, ok := line.Tokens[0].(token.LessToken); !ok {
		return false
	}

	return true
}

func checkIndentation(paragraph *token.TextParagraph) {
	var indent = -1
	for i, line := range paragraph.Lines {
		if indent == -1 || i == 1 {
			indent = int(line.Indentation)
		}

		if indent != int(line.Indentation) {
			return
		}

		paragraph.Indentation = indent
	}
}

func parseHeaderLines(paragraph token.TextParagraph) eNote.Options {
	res := eNote.NewOptions()

	for _, line := range paragraph.Lines {
		key, value := parseHeader(line.String())
		res.String[key] = value
	}

	return res
}

//isOnlyWhiteSpace returns true if the passed string contains only space, as defined by Unicode.
func isOnlyWhiteSpace(txt string) bool {
	txt = strings.TrimSpace(txt)
	return len(txt) == 0
}

//isCodeHeader checks if the line can be the header for a code block
//[langName]
func isCodeHeader(line token.LineContainer) bool {
	if len(line.Tokens) < 3 {
		return false
	}

	if _, ok := line.Tokens[0].(token.SBracketOpenToken); !ok {
		return false
	}
	if _, ok := line.Tokens[2].(token.SBracketCloseToken); !ok {
		return false
	}
	if text, ok := line.Tokens[1].(token.TextToken); !ok || len(text.Text) < 2 {
		return false
	}
	return true
}

func isQuoteLine(line token.LineContainer) bool {
	if len(line.Tokens) == 0 {
		return false
	}
	_, ok := line.Tokens[0].(token.PipeToken)
	return ok
}

//parseLine parses the passed token.LineContainer searching for special
func parseLine(currentLine token.LineContainer) token.LineToken {
	switch {
	case isType(token.TypeHeader, currentLine):
		return token.HeaderLine{}
	case isType(token.TypeEqual, currentLine, isTypeOptions{ignoreTabs: true}):
		log.Println("\t- Found EqualLine")
		return token.EqualLine{
			Indentation: currentLine.Indentation,
			Length:      len(currentLine.Tokens),
		}
	case isType(token.TypeLess, currentLine, isTypeOptions{ignoreTabs: true}):
		log.Println("\t- Found LessLine")
		return token.LessLine{
			Indentation: currentLine.Indentation,
			Length:      len(currentLine.Tokens),
		}
	case isListLine(currentLine):
		log.Println("\t- Found ListLine")
		indentation := 0
		for _, t := range currentLine.Tokens {
			if _, ok := t.(token.LessToken); ok {
				indentation++
			} else {
				currentLine.Tokens = currentLine.Tokens[indentation:]
				break
			}
		}

		fToken := currentLine.Tokens[0]
		if ft, ok := fToken.(token.TextToken); ok {
			ft.Text = ft.Text[1:]
			currentLine.Tokens[0] = ft
		}
		return token.ListLine{Text: currentLine, Indentation: indentation}
	case isCodeHeader(currentLine):
		log.Println("\t- Found CodeLine")
		lang := currentLine.Tokens[1].(token.TextToken).Text
		return token.CodeLine{Lang: lang}
	case isQuoteLine(currentLine):
		log.Println("\t- Found quote line")
		currentLine.Quote = true
		currentLine.Tokens = currentLine.Tokens[1:]
		return currentLine
	}

	return currentLine
}
