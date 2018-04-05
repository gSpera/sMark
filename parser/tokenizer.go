package parser

import (
	"eNote/token"
	"eNote/utils"
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
)

//Tokenizer parse a *os.File and return a slice of tokens
func Tokenizer(reader io.Reader) ([]token.Token, error) {
	tokenList := []token.Token{}
	char := make([]byte, 1)
	buffer := ""

	for {
		n, err := reader.Read(char)
		if n == 0 {
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.NewLineToken{})
			fmt.Println("EOF")
			tokenList = checkTokenList(tokenList)
			return tokenList, nil
		}
		if err != nil {
			return nil, err
		}

		switch char[0] {
		case token.TypeTab:
			fmt.Println("TAB")
			tokenList = append(tokenList, token.TabToken{})
		case token.TypeNewLine:
			fmt.Println("NewLine")
			if len(buffer) != 0 {
				addBufferToTokenBuffer(&tokenList, &buffer)
			}
			tokenList = append(tokenList, token.NewLineToken{})
		case token.TypeBold:
			fmt.Println("Bold")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.BoldToken{})

		case token.TypeItalic:
			fmt.Println("Italic")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.ItalicToken{})
		case token.TypeLess:
			fmt.Println("Less")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.LessToken{})
		case token.TypeHeader:
			fmt.Println("Header")
			if len(buffer) == 0 {
				// 	addBufferToTokenBuffer(&tokenList, &buffer)
				// } else {
				tokenList = append(tokenList, token.HeaderToken{})
				break
			}
		case token.TypeEqual:
			fmt.Println("Equal")
			addBufferToTokenBuffer(&tokenList, &buffer)
			tokenList = append(tokenList, token.EqualToken{})
		default:
			// fmt.Printf("Char: %c\n", char[0])
			buffer += string(char[0])
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
	lines := []token.LineToken{}
	currentLine := token.LineContainer{}

	indent := true

	isType := func(typ token.Type, line token.LineContainer) bool {
		fmt.Printf("isType: %c len: %v\n", typ, len(line.Tokens))
		if len(line.Tokens) < 2 {
			return false
		}

		for _, t := range line.Tokens {
			if typ != t.Type() {
				fmt.Printf("%c != %c\n", typ, t.Type())
				return false
			}
		}

		return true
	}

	for _, t := range tokens {
		switch t.(type) {
		case token.TabToken:
			if indent {
				fmt.Println("Adding Indentation")
				currentLine.Indentation++
			} else {
				indent = false
			}
		case token.NewLineToken:
			fmt.Println("NewLine")
			spew.Dump(currentLine.Tokens)

			switch {
			case isType(token.TypeHeader, currentLine):
				lines = append(lines, token.HeaderLine{})
				currentLine = token.LineContainer{}
				continue
			case isType(token.TypeEqual, currentLine):
				fmt.Println("EqualLine")
				lines = append(lines, token.EqualLine{})
				currentLine = token.LineContainer{}
				continue
			default:
				spew.Dump(currentLine.Tokens)
				fmt.Printf("TextLine: =: %v %T{%+v}\n", isType(token.TypeEqual, currentLine), currentLine, currentLine)
			}

			lines = append(lines, currentLine)
			currentLine = token.LineContainer{}
			continue
		case token.LessToken:
			fmt.Println("Less Token")
			if len(currentLine.Tokens) != 1 {
				continue //TODO: StrickeThrought
			}
			lines = append(lines, token.LessLine{})
			currentLine = token.LineContainer{}
			continue
		}

		currentLine.Tokens = append(currentLine.Tokens, t)
		indent = false
	}

	return lines
}

//TokenToParagraph divide a slice of lines in paragraphs
func TokenToParagraph(lines []token.LineToken) []token.ParagraphToken {
	fmt.Printf("Paragraphs: %d Lines\n", len(lines))
	paragraphs := []token.ParagraphToken{}
	currentParagraph := token.TextParagraph{}
	header := false

	for i, t := range lines {
		var lastLine token.LineToken
		if i == 0 {
			lastLine = nil
		} else {
			lastLine = lines[i-1]
		}
		switch t.(type) {
		case token.EqualLine:
			fmt.Println("EqualLine Paragraph")
			fmt.Println(len(currentParagraph.Lines))
			if len(currentParagraph.Lines) != 2 {
				continue
			}

			if _, ok := lastLine.(token.LineContainer); !ok {
				continue
			}
			lastLine := lastLine.(token.LineContainer)

			if len(lastLine.Tokens) == 0 {
				continue
			}
			fmt.Println("Title Paragraph")
			paragraphs = append(paragraphs, token.TitleParagraph{Text: lastLine})
			currentParagraph = token.TextParagraph{}
		case token.HeaderLine:
			fmt.Println("HeaderLine")
			if !header && len(currentParagraph.Lines) != 0 {
				fmt.Println("Not at first")
				continue
			}

			if header {
				fmt.Println("Appending Header Paragraph")
				paragraphs = append(paragraphs, token.HeaderParagraph{OptionsTemplate: parseHeaderLines(currentParagraph)})
				currentParagraph = token.TextParagraph{}
			}

			header = !header
		case token.LessLine:
			fmt.Println("LessLine")
			//Subtitle or divisor
			if _, ok := lastLine.(token.LineContainer); !ok {
				continue
			}
			switch len(lastLine.(token.LineContainer).Tokens) {
			case 0:
				fmt.Println("Divisor")
				paragraphs = append(paragraphs, token.DivisorParagraph{})
				currentParagraph = token.TextParagraph{}
			default: //TODO: Subtitle
				fmt.Println("Subtitle (NOT IMPLEMENTED)")
			}

		case token.LineContainer:
			fmt.Println("LineContainer")
			t := t.(token.LineContainer)
			currentEmpty := len(t.Tokens) == 0
			spew.Dump(t.Tokens)

			if l, ok := lastLine.(token.LineContainer); currentEmpty && ok && len(l.Tokens) != 0 {
				fmt.Println("New Paragraph")
				fmt.Printf("LastLine: Len: %d %T: %+v\n", len(l.Tokens), l, l)
				spew.Dump(l.Tokens)
				checkIndentation(&currentParagraph)
				fmt.Printf("Indentation after: %d\n", currentParagraph.Indentation)
				paragraphs = append(paragraphs, currentParagraph)
				currentParagraph = token.TextParagraph{}
			} else {
				fmt.Println("Old Paragraph")
				spew.Dump(t.Tokens)
				currentParagraph.Lines = append(currentParagraph.Lines, t)
			}
		}

	}

	checkIndentation(&currentParagraph)
	paragraphs = append(paragraphs, currentParagraph)
	return paragraphs
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

func parseHeaderLines(paragraph token.TextParagraph) eNote.OptionsTemplate {
	res := eNote.OptionsTemplate{}

	for _, line := range paragraph.Lines {
		key, value := parseHeader(line.String())
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		res.AddString(key, value)
	}

	return res
}
