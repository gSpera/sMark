package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	output "eNote/output/html"
	"eNote/output/telegraph"
	"eNote/parser"
	eNote "eNote/utils"

	"github.com/davecgh/go-spew/spew"
	tgraph "github.com/toby3d/telegraph"
)

//ProgramName is the name of the program.
const ProgramName = "eNote"

func main() {
	//Logger
	log.SetPrefix(ProgramName + ": ")
	log.SetFlags(log.Ltime | log.Lshortfile)
	verbose := flag.Bool("verbose", false, "Output logging")
	//Flags
	options := eNote.Options{
		InputFile:  flag.String("i", "-", "Input file, - for stdin"),
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
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	input, err := streamFromFilename(*options.InputFile)
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Error: File doesn't exist")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("Error: Could not open file: %v", err)
	}

	fmt.Println("Filename: ", *options.InputFile)
	fmt.Println("Input: ", input)

	reader := bufio.NewReader(input)

	log.Println("Resetting Reader")
	input.Seek(0, os.SEEK_SET)
	reader = bufio.NewReader(input)

	log.Println("Parsing")
	tokenList, err := parser.ParseReader(reader)
	if err != nil {
		log.Fatalf("Error: Could not parse: %v", err)
	}
	log.Println("Parsing DONE")
	log.Println("Updating Options")
	options.Update(parser.OptionsFromParagraphs(&tokenList))
	log.Println("Updating Options DONE")

	spew.Dump(tokenList)
	fmt.Printf("%v\n", len(tokenList))

	fmt.Println("Final Options")
	spew.Dump(options)

	//TODO: Add more flexibility to the chose of output engine
	//HTML Output Engine
	log.Println("Outputting HTML")
	ioutil.WriteFile(*options.OutputFile, output.ToString(tokenList, options), os.ModePerm)
	log.Println("Outputting HTML DONE")

	//Telegraph Output Engine
	log.Println("Outputting Telegraph")
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
		log.Fatalf("Error: Could not create page: %v", err)
	}

	spew.Dump(pagePubblished)
	log.Println("Outputting Telegraph DONE")

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
