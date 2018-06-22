package parser

import (
	"bufio"
	"eNote/token"
	"eNote/utils"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

//Tokenizer parse a *os.File and return a slice of tokens
func Tokenizer(reader io.Reader) ([]token.Token, error) {
	tokenList := []token.Token{}
	r := bufio.NewReader(reader)
	buffer := ""

	for {
		// n, err := reader.Read(char)
		n, size, err := r.ReadRune()

		if size == 0 {
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.NewLineToken{})
			fmt.Println("EOF")
			tokenList = checkTokenList(tokenList)
			return tokenList, nil
		}
		if err != nil {
			return nil, err
		}

		switch n {
		case token.TypeTab:
			log.Println("\t- Found TabToken")
			tokenList = append(tokenList, token.TabToken{})
		case token.TypeNewLine:
			log.Println("\t- Found NewLineToken")
			if len(buffer) != 0 {
				addBufferToTokenBuffer(&tokenList, &buffer)
			}
			tokenList = append(tokenList, token.NewLineToken{})
		case token.TypeBold:
			log.Println("\t- Found BoldToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.BoldToken{})
		case token.TypeItalic:
			log.Println("\t- Found ItalicToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.ItalicToken{})
		case token.TypeLess:
			log.Println("\t- Found LessToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.LessToken{})
		case token.TypeHeader:
			log.Println("\t- Found HeaderToken")
			if len(buffer) == 0 {
				// 	addBufferToTokenBuffer(&tokenList, &buffer)
				// } else {
				tokenList = append(tokenList, token.HeaderToken{})
				break
			}
		case token.TypeEqual:
			log.Println("\t- Found EqualToken")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.EqualToken{})
		default:
			// fmt.Printf("Char: %c\n", char[0])
			buffer += string(n)
		}
	}
}

func addBufferToTokenBuffer(tokenBuffer *[]token.Token, buffer *string) {
	if len(*buffer) == 0 {
		return
	}

	*tokenBuffer = append(*tokenBuffer, token.TextToken{Text: *buffer})
	*buffer = ""
}

func checkTokenList(tokenList []token.Token) []token.Token {
	if len(tokenList) == 0 {
		return tokenList
	}

	if t, ok := tokenList[0].(token.TextToken); ok && len(t.Text) == 0 {
		return tokenList[1:]
	}

	return tokenList
}

//TokenToLine divide a slice of tokens in lines
func TokenToLine(tokens []token.Token) []token.LineToken {
	fmt.Println("TokenToLine")
	lines := []token.LineToken{}
	currentLine := token.LineContainer{}

	indent := true

	for _, t := range tokens {
		switch t.(type) {
		case token.TabToken:
			if indent {
				fmt.Println("Adding Indentation")
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
	fmt.Printf("isType: %c len: %v\n", typ, len(line.Tokens))
	fmt.Println(line)

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
			fmt.Println("Ignoring Tab")
			continue
		}
		if typ != t.Type() {
			fmt.Printf("%c != %c\n", typ, t.Type())
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
	fmt.Println("Check indentation")
	for i, line := range paragraph.Lines {
		fmt.Printf("Line %d Indentation: %d\n", i, line.Indentation)
		if indent == -1 || i == 1 {
			indent = int(line.Indentation)
		}

		if indent != int(line.Indentation) {
			fmt.Printf("Check indetation %d != %d\n", indent, line.Indentation)
			return
		}

		paragraph.Indentation = indent
		fmt.Printf("New Indetation: %d\n", indent)
	}
}

func parseHeaderLines(paragraph token.TextParagraph) eNote.Options {
	res := eNote.Options{}

	for _, line := range paragraph.Lines {
		key, value := parseHeader(line.String())
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		res.AddString(key, value)
	}

	return res
}

//isOnlyWhiteSpace returns true if the passed string contains only space, as defined by Unicode.
func isOnlyWhiteSpace(txt string) bool {
	txt = strings.TrimSpace(txt)
	return len(txt) == 0
}

//parseLine parses the passed token.LineContainer searching for special
func parseLine(currentLine token.LineContainer) token.LineToken {
	fmt.Println("NewLine")
	spew.Dump(currentLine.Tokens)

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
		currentLine.Tokens = currentLine.Tokens[1:]
		fToken := currentLine.Tokens[0]
		if ft, ok := fToken.(token.TextToken); ok {
			ft.Text = ft.Text[1:]
			currentLine.Tokens[0] = ft
		}
		return token.ListLine{Text: currentLine}
	default:
		spew.Dump(currentLine.Tokens)
		fmt.Printf("TextLine: =: %v %T{%+v}\n", isType(token.TypeEqual, currentLine), currentLine, currentLine)
	}

	return currentLine
}
