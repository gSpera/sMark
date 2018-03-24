package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	filename := flag.String("i", "-", "input file, - for stdin, default is stdin")
	flag.Parse()
	input, err := streamFromFilename(*filename)
	if os.IsNotExist(err) {
		fmt.Println("Error: File doesn't exist")
		printUsage()
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("Filename: ", *filename)
	fmt.Println("Input: ", input)
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
