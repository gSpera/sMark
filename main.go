package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"

	"eNote/output"
	"eNote/parser"
)

func main() {
	inputfile := flag.String("i", "-", "input file, - for stdin")
	outfile := flag.String("o", "out.rtf", "output file")

	flag.Parse()
	input, err := streamFromFilename(*inputfile)
	if os.IsNotExist(err) {
		fmt.Println("Error: File doesn't exist")
		printUsage()
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("Filename: ", *inputfile)
	fmt.Println("Input: ", input)
	tokenList, err := parser.ParseFile(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TokenList %d tokens: %+v\n", len(tokenList), tokenList)
	for _, t := range tokenList {
		fmt.Printf("%s\n", t.DebugString())
	}
	spew.Dump(tokenList)
	ioutil.WriteFile(*outfile, []byte(output.ToString(tokenList)), os.ModePerm)
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

func printUsage() {
	fmt.Printf("%s -i <input file>\n", os.Args[0])
}
