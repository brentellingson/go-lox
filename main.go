package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/brentellingson/go-lox/internal/errutil"
	"github.com/brentellingson/go-lox/internal/scanner"
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
	if errutil.HadError {
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
		errutil.HadError = false
	}

	if err := scanner.Err(); err != nil {
		panic("error reading input")
	}
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	for _, token := range scanner.ScanTokens() {
		fmt.Println(token)
	}
}
