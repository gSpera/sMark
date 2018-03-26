package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"

	"eNote/output"
	"eNote/parser"
	eNote "eNote/utils"
)

func main() {
	options := eNote.Options{
		InputFile:  flag.String("i", "-", "IUnput file, - for stdin"),
		OutputFile: flag.String("o", "out.html", "Output file"),
		NewLine:    flag.Bool("newline", true, "Include newlines as in source"),
		CustomCSS:  flag.String("css", "", "A custom css file"),
		InlineCSS:  flag.String("inline-css", "", "Inline CSS"),
		EnableFont: flag.Bool("font", true, "Enable a default font"),
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
	tokenList, err := parser.ParseReader(input)
	if err != nil {
		panic(err)
	}

	spew.Dump(tokenList)
	fmt.Printf("%v\n", len(tokenList))
	ioutil.WriteFile(*options.OutputFile, output.ToString(tokenList, options), os.ModePerm)
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
