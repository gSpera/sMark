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
		log.Println("Rune:", n, string(n), "Size:", size, "Err:", err)
		if size != 1 {
			log.Println()
		}
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
			fmt.Println("NewLine")
			spew.Dump(currentLine.Tokens)

			switch {
			case isType(token.TypeHeader, currentLine):
				lines = append(lines, token.HeaderLine{})
				currentLine = token.LineContainer{}
				continue
			case isType(token.TypeEqual, currentLine, isTypeOptions{ignoreTabs: true}):
				log.Println("\t- Found EqualLine")
				log.Println("\t\t- Indentation:", currentLine.Indentation)

				lines = append(lines, token.EqualLine{
					Indentation: currentLine.Indentation,
					Length:      uint(len(currentLine.Tokens)),
				})

				currentLine = token.LineContainer{}
				continue
			case isType(token.TypeLess, currentLine, isTypeOptions{ignoreTabs: true}):
				log.Println("\t- Found LessLine")
				lines = append(lines, token.LessLine{})
				currentLine = token.LineContainer{}
				continue
			case isListLine(currentLine):
				log.Println("\t- Found ListLine")
				currentLine.Tokens = currentLine.Tokens[1:]
				fToken := currentLine.Tokens[0]
				if ft, ok := fToken.(token.TextToken); ok {
					ft.Text = ft.Text[1:]
					currentLine.Tokens[0] = ft
				}
				lines = append(lines, token.ListLine{Text: currentLine})
				currentLine = token.LineContainer{}
				continue
			default:
				log.Println("\t- Found TextLine")
				log.Println("\t\t- Indentation:", currentLine.Indentation)
				spew.Dump(currentLine.Tokens)
				fmt.Printf("TextLine: =: %v %T{%+v}\n", isType(token.TypeEqual, currentLine), currentLine, currentLine)
			}

			lines = append(lines, currentLine)
			currentLine = token.LineContainer{}
			indent = true
			continue
			// case token.LessToken:
			// 	fmt.Println("Less Token")
			// 	if !isType(token.TypeLess, currentLine, isTypeOptions{threshold: 0}) {
			// 		//Strick-throught
			// 		panic("Strickthrought")
			// 		continue
			// 	}
		}

		currentLine.Tokens = append(currentLine.Tokens, t)
		indent = false
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

//TokenToParagraph divide a slice of lines in paragraphs
func TokenToParagraph(lines []token.LineToken) []token.ParagraphToken {
	fmt.Printf("Paragraphs: %d Lines\n", len(lines))
	paragraphs := []token.ParagraphToken{}
	currentParagraph := token.TextParagraph{}
	header := false
	var list []token.ListLine

	for i, t := range lines {
		var lastLine token.LineToken
		if i == 0 {
			lastLine = nil
		} else {
			lastLine = lines[i-1]
		}
		if _, ok := t.(token.ListLine); !ok && list != nil {
			paragraphs = append(paragraphs, token.ListParagraph{Items: list})
			list = nil
		}

		switch tt := t.(type) {
		case token.EqualLine:
			log.Println("\t- Found EqualLine Line")
			fmt.Println(len(currentParagraph.Lines))

			if notEmptyLines(currentParagraph.Lines) != 1 {
				// log.Println(currentParagraph)
				log.Println("\t\t- Wrong number of lines:", notEmptyLines(currentParagraph.Lines))
				currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("=", int(tt.Length))))
				continue
			}

			if _, ok := lastLine.(token.LineContainer); !ok {
				log.Println("\t\t- LastLine is not a token.LineContainer")
				currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("=", int(tt.Length))))
				continue
			}
			lastLine := lastLine.(token.LineContainer)

			if len(lastLine.Tokens) == 0 {
				log.Println("\t\t- LastLine is empty")
				currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("=", int(tt.Length))))
				continue
			}

			if lastLine.Indentation != uint(lastLine.Indentation) {
				log.Println("\t\t- Indentation are differents", lastLine.Indentation, lastLine.Indentation)
				currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("=", int(tt.Length))))
				continue
			}

			log.Println("\t- Found Title Paragraph")
			paragraphs = append(paragraphs, token.TitleParagraph{Text: lastLine, Indentation: lastLine.Indentation})
			currentParagraph = token.TextParagraph{}
		case token.HeaderLine:
			fmt.Println("HeaderLine")
			if !header && len(currentParagraph.Lines) != 0 {
				fmt.Println("Not at first")
				continue
			}

			if header {
				log.Println("\t- Found Header Paragraph")
				paragraphs = append(paragraphs, token.HeaderParagraph{Options: parseHeaderLines(currentParagraph)})
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
				log.Println("\t- Found Divisor")
				paragraphs = append(paragraphs, token.DivisorParagraph{})
				currentParagraph = token.TextParagraph{}
			default:
				if notEmptyLines(currentParagraph.Lines) != 1 {
					// log.Println(currentParagraph)
					log.Println("\t\t- Wrong number of lines:", notEmptyLines(currentParagraph.Lines))
					currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("-", int(tt.Length))))
					continue
				}

				if _, ok := lastLine.(token.LineContainer); !ok {
					log.Println("\t\t- LastLine is not a token.LineContainer")
					currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("-", int(tt.Length))))
					continue
				}
				lastLine := lastLine.(token.LineContainer)

				if len(lastLine.Tokens) == 0 {
					log.Println("\t\t- LastLine is empty")
					currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("-", int(tt.Length))))
					continue
				}

				if lastLine.Indentation != uint(lastLine.Indentation) {
					log.Println("\t\t- Indentation are differents", lastLine.Indentation, lastLine.Indentation)
					currentParagraph.Lines = append(currentParagraph.Lines, token.LineContainerFromString(strings.Repeat("-", int(tt.Length))))
					continue
				}

				paragraphs = append(paragraphs, token.SubtitleParagraph{
					Text:        lastLine,
					Indentation: lastLine.Indentation,
				})
				currentParagraph = token.TextParagraph{}
			}

		case token.LineContainer:
			fmt.Println("LineContainer")
			t := t.(token.LineContainer)
			currentEmpty := len(t.Tokens) == 0
			spew.Dump(t.Tokens)

			if l, ok := lastLine.(token.LineContainer); currentEmpty && ok && len(l.Tokens) != 0 {
				log.Println("\t- Found Text Paragraph")
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
		case token.ListLine:
			list = append(list, tt)
		default:
			panic(fmt.Sprintf("Line=>Paragraph for %T{%+v} not defined", tt, tt))
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

func parseHeaderLines(paragraph token.TextParagraph) eNote.Options {
	res := eNote.Options{}

	for _, line := range paragraph.Lines {
		key, value := parseHeader(line.String())
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		res.AddString(key, value)
	}

	return res
}

func notEmptyLines(lines []token.LineContainer) uint {
	notEmpty := uint(0)

	for _, l := range lines {
		if len(l.Tokens) == 0 {
			continue
		}

		for _, t := range l.Tokens {
			switch tt := t.(type) {
			case token.TabToken:
				continue
			case token.TextToken:
				if isOnlyWhiteSpace(tt.Text) {
					continue
				}
			}
		}

		notEmpty++
	}

	return notEmpty
}

//isOnlyWhiteSpace returns true if the passed string contains only space, as defined by Unicode.
func isOnlyWhiteSpace(txt string) bool {
	txt = strings.TrimSpace(txt)
	return len(txt) == 0
}
