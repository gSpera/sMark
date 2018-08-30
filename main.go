package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gSpera/sMark/output/html"
	"github.com/gSpera/sMark/output/prettify"
	"github.com/gSpera/sMark/output/telegraph"
	"github.com/gSpera/sMark/parser"
	"github.com/gSpera/sMark/token"
	sMark "github.com/gSpera/sMark/utils"

	tgraph "github.com/toby3d/telegraph"
	fsnotify "gopkg.in/fsnotify.v1"
)

//ProgramName is the name of the program.
const ProgramName = "sMark"

func main() {
	//Logger
	log.SetPrefix(ProgramName + ": ")
	log.SetFlags(log.Ltime | log.Lshortfile)

	//Options from Command Line
	inputFile := flag.String("i", "-", "Input file, - for stdin")
	outputFile := flag.String("o", "out.html", "Output file")
	newLine := flag.Bool("newline", true, "Include newlines as in source")
	customCSS := flag.String("css", "", "A custom css file")
	inlineCSS := flag.String("inline-css", "", "Inline CSS")
	enableFont := flag.Bool("font", true, "Enable a default font")
	onlyBody := flag.Bool("only-body", false, "Output only the html boy and not the whole page")
	title := flag.String("title", "", "The title of the output document")
	tabWidth := flag.Uint("tabs-width", 4, "The width (in spaces) of one tab")
	verbose := flag.Bool("verbose", false, "Output logging")
	htmlOut := flag.Bool("html", true, "Output HTML")
	telegraphOut := flag.Bool("telegraph", false, "Output to Telegra.ph")
	prettify := flag.String("prettify", "", "Prettify")
	watch := flag.Bool("watch", false, "Watches the file for changes (CTRL + c to exit)")
	enableJs := flag.Bool("enable-js", true, "Enable JavaScript in HTML")

	flag.Parse()
	options := sMark.Options{
		String: map[string]string{
			"InputFile":  *inputFile,
			"OutputFile": *outputFile,
			"CustomCSS":  *customCSS,
			"InlineCSS":  *inlineCSS,
			"Title":      *title,
			"Prettify":   *prettify,
		},
		Bool: map[string]bool{
			"NewLine":      *newLine,
			"EnableFont":   *enableFont,
			"OnlyBody":     *onlyBody,
			"Verbose":      *verbose,
			"HTMLOut":      *htmlOut,
			"TelegraphOut": *telegraphOut,
			"Watch":        *watch,
			"EnableJS":     *enableJs,
		},
		Generic: map[string]interface{}{
			"TabWidth": *tabWidth, //Uint
		},
	}

	if !options.Bool["Verbose"] {
		log.SetOutput(ioutil.Discard)
	}

	if options.Bool["Watch"] {
		if options.String["InputFile"] == "-" {
			fmt.Fprintln(os.Stderr, "Could not watch stdin")
			return
		}
		fmt.Println("Watching File, use CTRL + c to stop")
		options, _ = compile(options)

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create watcher: %v", err)
		}
		defer watcher.Close()

		err = watcher.Add(options.String["InputFile"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot watch %s: %v\n", options.String["InputFile"], err)
		}

		if tmpl, ok := options.String["TemplateFile"]; ok {
			if err := watcher.Add(options.String["TemplateFile"]); err != nil {
				fmt.Fprintf(os.Stderr, "Could not Watch TemplateFile: %s: %v\n", tmpl, err)
			} else {
				fmt.Fprintf(os.Stderr, "Watching TemplateFile: %s\n", tmpl)
			}
		}

		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
					fmt.Fprintln(os.Stderr, "Quitting, file is being removed or renamed")
					return
				}
				if event.Op != fsnotify.Write {
					continue
				}
				compile(options)
			case err := <-watcher.Errors:
				fmt.Fprintln(os.Stderr, "Error in watcher:", err)
			}
		}
	}
	compile(options)
}

func compile(options sMark.Options) (sMark.Options, []token.ParagraphToken) {
	log.Println("Start Compilation")
	input, err := streamFromFilename(options.String["InputFile"])
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Error: File doesn't exist")
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("Error: Could not open file: %v", err)
	}

	log.Println("Filename: ", options.String["InputFile"])

	reader := bufio.NewReader(input)

	log.Println("Parsing")
	tokenList, err := parser.ParseReader(reader)
	if err != nil {
		log.Fatalf("Error: Could not parse: %v", err)
	}
	log.Println("Parsing DONE")

	log.Println("Updating Options")
	log.Println("\t- From Paragraphs")
	newOptions, tokenList := parser.OptionsFromParagraphs(tokenList)
	flOptions := newOptions
	options.Update(newOptions)
	log.Println("\t- From Paragraphs DONE")
	log.Println("\t- From Command Line")
	for _, arg := range flag.Args() {
		k, v := sMark.ParseOption(arg)
		options.Insert(k, v)
	}
	log.Println("\t- From Command Line DONE")
	log.Println("Updating Options DONE")

	log.Println("Selecting Title")
	if options.String["Title"] == "" {
		log.Println("Changing Title")
		title := parser.TitleFromParagraph(tokenList)
		options.String["Title"] = title
	}
	log.Println("Selecting Title DONE")

	//HTML Output Engine
	if options.Bool["HTMLOut"] {
		log.Println("Outputting HTML")
		ioutil.WriteFile(options.String["OutputFile"], htmlout.ToString(tokenList, options), os.ModePerm)
		log.Println("Outputting HTML DONE")
	}

	//Telegraph Output Engine
	if options.Bool["TelegraphOut"] {
		log.Println("Outputting Telegraph")
		page := telegraphout.ToString(tokenList, options)
		log.Println("Title:", page.Title)
		var accessToken string
		fmt.Print("Insert Access Token: ")
		fmt.Scan(&accessToken)
		acc := tgraph.Account{
			AccessToken: accessToken,
		}
		pagePubblished, err := acc.CreatePage(&page, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot create page:", err)
			log.Fatalf("Error: Could not create page: %v", err)
		}

		log.Println("Outputting Telegraph DONE")
		fmt.Println("Title:", pagePubblished.Title)
		fmt.Println("URL:", pagePubblished.URL)
	}

	//sMark Output Engine
	if options.String["Prettify"] != "" {
		log.Println("Outputting sMark")
		//Prettifier doesn't need all the Options but only from the file
		data, err := prettify.Output(tokenList, flOptions)
		if err != nil {
			log.Fatalf("Could not compile to sMark: %v\n", err)
		}

		f, err := os.Create(options.String["Prettify"])
		if err != nil {
			log.Fatalf("Could not create file %s: %v", options.String["Prettify"], err)
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			log.Fatalf("Could not output to file: %v\n", err)
		}

		log.Println("Outputting sMark DONE")
	}

	return options, tokenList
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
