package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/brentellingson/go-lox/internal"
	"github.com/brentellingson/go-lox/internal/parse"
	"github.com/brentellingson/go-lox/internal/scan"
	"github.com/brentellingson/go-lox/internal/vm"
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
	if internal.HadRuntimeError {
		os.Exit(70)
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
		internal.HadRuntimeError = false
	}

	if err := scanner.Err(); err != nil {
		panic("error reading input")
	}
}

func run(source string) {
	scanner := scan.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := parse.NewParser(tokens)
	expr := parser.Parse()

	if internal.HadParserError {
		return
	}

	interpreter := vm.Interpreter{}
	interpreter.Interpret(expr)
}
