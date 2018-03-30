package parser

import (
	"bufio"
	"eNote/token"
	"eNote/utils"
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
)

//ParseReader parse a source file and returns the paragraphs
func ParseReader(fl io.Reader) ([]token.ParagraphToken, error) {
	tokens, err := Tokenizer(fl)
	spew.Dump(tokens)
	if err != nil {
		return nil, err
	}

	lines := TokenToLine(tokens)
	paragraphs := TokenToParagraph(lines)

	return paragraphs, nil
}

//ParseHeader parses the header from the header, returns the eNote.Options and a bool rappresenting the status
func ParseHeader(r *bufio.Reader) (eNote.OptionsTemplate, bool) {
	res := eNote.OptionsTemplate{}

	line, err := r.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		return res, false
	}

	//Check if line is a header starting
	for _, ch := range line[:1] {
		if ch != '+' {
			fmt.Println("First Line is not heading start")
			fmt.Println(line)
			return res, false
		}
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		//Check if line is a header ending
		ending := true
		for _, ch := range line[:1] {
			if ch != '+' {
				ending = false
			}
		}
		if ending {
			return res, true
		}

		key, value := parseHeader(line)
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		res.AddString(key, value)
	}

	return res, true
}

func parseHeader(line string) (string, string) {
	key := ""
	buffer := []rune{}
	fmt.Println(len(line))
	for _, ch := range line {
		switch ch {
		case '=':
			key = string(buffer)
			buffer = []rune{}
		default:
			buffer = append(buffer, ch)
		}
	}

	return key, string(buffer)
}
