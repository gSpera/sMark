package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	output "eNote/output/html"
	"eNote/output/telegraph"
	"eNote/parser"
	eNote "eNote/utils"

	"github.com/davecgh/go-spew/spew"
	tgraph "github.com/toby3d/telegraph"
)

func main() {
	options := eNote.Options{
		InputFile:  flag.String("i", "-", "IUnput file, - for stdin"),
		OutputFile: flag.String("o", "out.html", "Output file"),
		NewLine:    flag.Bool("newline", true, "Include newlines as in source"),
		CustomCSS:  flag.String("css", "", "A custom css file"),
		InlineCSS:  flag.String("inline-css", "", "Inline CSS"),
		EnableFont: flag.Bool("font", true, "Enable a default font"),
		OnlyBody:   flag.Bool("only-body", false, "Output only the html boy and not the whole page"),
		Title:      flag.String("title", "", "The title of the output document"),
		TabWidth:   flag.Uint("tabs-width", 4, "The width (in spaces) of one tab"),
	}

	flag.Parse()
	input, err := streamFromFilename(*options.InputFile)
	if os.IsNotExist(err) {
		fmt.Println("Error: File doesn't exist")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("Filename: ", *options.InputFile)
	fmt.Println("Input: ", input)

	reader := bufio.NewReader(input)
	header, ok := eNote.OptionsTemplate{}, false // parser.ParseHeader(reader)
	spew.Dump(header)

	if ok {
		fmt.Println("Updating")
		options.Update(header)
		spew.Dump(options)
	} else {
		fmt.Println("Resetting Reader")
		input.Seek(0, os.SEEK_SET)
		reader = bufio.NewReader(input)
	}
	tokenList, err := parser.ParseReader(reader)
	if err != nil {
		panic(err)
	}

	spew.Dump(tokenList)
	fmt.Printf("%v\n", len(tokenList))
	//TODO: Add more flexibility to the chose of output engine
	//HTML Output Engine
	ioutil.WriteFile(*options.OutputFile, output.ToString(tokenList, options), os.ModePerm)

	//Telegraph Output Engine
	page := outTelegraph.ToString(tokenList, options)
	fmt.Println("Title:", page.Title)
	spew.Dump(page)
	var accessToken string
	fmt.Println("Insert Access Token: ")
	fmt.Scan(&accessToken)
	acc := tgraph.Account{
		AccessToken: accessToken,
	}
	pagePubblished, err := acc.CreatePage(&page, false)
	if err != nil {
		fmt.Println("Error: Could not create page:")
		fmt.Println(err.Error())
	}
	spew.Dump(pagePubblished)
}

func streamFromFilename(filename string) (*os.File, error) {
	if filename == "-" {
		return os.Stdin, nil
	}

	fl, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return fl, nil
}
