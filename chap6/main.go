package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/brentellingson/go-lox/internal"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: go-lox [script]")
		os.Exit(64)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic("error reading file " + path)
	}
	run(string(bytes))
	if internal.HadParserError {
		os.Exit(65)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line)
		internal.HadParserError = false
	}

	if err := scanner.Err(); err != nil {
		panic("error reading input")
	}
}

func run(source string) {
	scanner := internal.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := internal.NewParser(tokens)
	expr := parser.Parse()

	if internal.HadParserError {
		return
	}

	fmt.Println(internal.PrintAst(expr))
}
