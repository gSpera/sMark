package parser

import (
	"bufio"
	"eNote/token"
	"eNote/utils"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/davecgh/go-spew/spew"
)

//ParseReader parse a source file and returns the paragraphs
func ParseReader(fl io.Reader) ([]token.ParagraphToken, error) {
	log.Println("Tokenizer")
	tokens, err := Tokenizer(fl)
	log.Println("Tokenizer DONE")
	spew.Dump(tokens)
	if err != nil {
		return nil, err
	}

	log.Println("Searching Structures")
	tokens = TokenToStructure(tokens)
	log.Println("Searching Structures")

	log.Println("Token To Line")
	lines := TokenToLine(tokens)
	log.Println("Token To Line DONE")
	ioutil.WriteFile("lines.dump", []byte(spew.Sdump(lines)), 0665)

	log.Println("Line To Paragraph")
	paragraphs := TokenToParagraph(lines)
	log.Println("Line To Paragraph DONE")

	spew.Dump(paragraphs)
	return paragraphs, nil
}

//ParseHeader parses the header from the header, returns the eNote.Options and a bool rappresenting the status
func ParseHeader(r *bufio.Reader) (eNote.Options, bool) {
	res := eNote.Options{}

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

//OptionsFromParagraphs analyze the passed slice of paragraphs returning the final options contined in token.HeaderParagraphs.
//It update the slice removing any HeaderParagraph token
func OptionsFromParagraphs(paragraphs *[]token.ParagraphToken) eNote.Options {
	fmt.Println("OptionsFromParagraphs")
	options := eNote.Options{}

	for i, p := range *paragraphs {
		fmt.Println("Paragraph")
		if _, ok := p.(token.HeaderParagraph); !ok {
			continue
		}
		fmt.Println(" - Header")
		p := p.(token.HeaderParagraph)
		spew.Dump(p)
		options.Update(p.Options)

		fmt.Println("After Update")
		spew.Dump(options)

		//Removing the paragraph
		par := *paragraphs
		*paragraphs = append(par[:i], par[i+1:]...)
	}

	return options
}

//TitleFromParagraph resolve a possible title from the passed paragraphs
func TitleFromParagraph(paragraph []token.ParagraphToken) string {
	fmt.Println("TitleFromParagraph")
	titles := map[string]int{}

	for _, p := range paragraph {
		if pp, ok := p.(token.TitleParagraph); ok {
			txt := pp.Text
			value := pp.Indentation

			if currentValue, ok := titles[txt]; !ok || currentValue < value {
				titles[txt] = value
			}
		}

	}

	max := "No Title"
	maxV := 0

	for k, v := range titles {
		if v > maxV {
			max = k
			maxV = v
		}
	}

	return max
}
