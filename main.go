package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	output "eNote/output/html"
	"eNote/output/prettify"
	"eNote/output/telegraph"
	"eNote/parser"
	eNote "eNote/utils"

	tgraph "github.com/toby3d/telegraph"
	fsnotify "gopkg.in/fsnotify.v1"
)

//ProgramName is the name of the program.
const ProgramName = "eNote"

func main() {
	//Logger
	log.SetPrefix(ProgramName + ": ")
	log.SetFlags(log.Ltime | log.Lshortfile)
	verbose := flag.Bool("verbose", false, "Output logging")
	htmlOut := flag.Bool("html", true, "Output HTML")
	telegraphOut := flag.Bool("telegraph", false, "Output to Telegra.ph")
	prettifyFlag := flag.String("prettify", "", "Prettify")
	watchFlag := flag.Bool("watch", false, "Watches the file for changes (CTRL + c to exit)")

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

	if *watchFlag {
		fmt.Println("Watching File, use CTRL + c to stop")
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create watcher: %v", err)
		}
		defer watcher.Close()

		err = watcher.Add(*options.InputFile)
		if err != nil {
			fmt.Printf("Cannot watch %s: %v\n", *options.InputFile, err)
		}

		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
					fmt.Println("Quitting, file is being removed or renamed")
					return
				}
				if event.Op != fsnotify.Write {
					continue
				}
				compile(options, verbose, htmlOut, telegraphOut, prettifyFlag)
			case err := <-watcher.Errors:
				fmt.Println("Error in watcher:", err)
			}
		}
	}

	compile(options, verbose, htmlOut, telegraphOut, prettifyFlag)
}

func compile(options eNote.Options, verbose, htmlOut, telegraphOut *bool, prettifyFlag *string) {
	input, err := streamFromFilename(*options.InputFile)
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Error: File doesn't exist")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("Error: Could not open file: %v", err)
	}

	log.Println("Filename: ", *options.InputFile)

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

	log.Println("Selecting Title")
	if options.Title == nil || *options.Title == "" {
		log.Println("Changing Title")
		title := parser.TitleFromParagraph(tokenList)
		options.Title = &title
	}
	log.Println("Selecting Title DONE")

	//HTML Output Engine
	if *htmlOut {
		log.Println("Outputting HTML")
		ioutil.WriteFile(*options.OutputFile, output.ToString(tokenList, options), os.ModePerm)
		log.Println("Outputting HTML DONE")
	}

	//Telegraph Output Engine
	if *telegraphOut {
		log.Println("Outputting Telegraph")
		page := outTelegraph.ToString(tokenList, options)
		log.Println("Title:", page.Title)
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

		log.Println("Outputting Telegraph DONE")
		fmt.Println("Title:", pagePubblished.Title)
		fmt.Println("URL:", pagePubblished.URL)
	}

	//eNote Output Engine
	if *prettifyFlag != "" {
		log.Println("Outputting eNote")
		data, err := prettify.Output(tokenList, options)
		if err != nil {
			log.Fatalf("Could not compile to eNote: %v\n", err)
		}

		f, err := os.Create(*prettifyFlag)
		if err != nil {
			log.Fatalf("Could not create file %s: %v", *prettifyFlag, err)
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			log.Fatalf("Could not output to file: %v\n", err)
		}

		log.Println("Outputting eNote DONE")
	}
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
