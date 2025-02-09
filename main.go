package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/brentellingson/go-lox/internal/engine"
	"github.com/brentellingson/go-lox/internal/parse"
	"github.com/brentellingson/go-lox/internal/repl"
	"github.com/brentellingson/go-lox/internal/scan"
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
	repl := repl.NewRepl(scan.Scan, parse.Parse, &engine.Interpreter{})
	_, err = repl.Run(string(bytes))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runPrompt() {
	repl := repl.NewRepl(scan.Scan, parse.Parse, engine.NewInterpreter())
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		rslt, err := repl.Run(line)
		if err != nil {
			fmt.Println(err)
		} else if rslt != nil {
			fmt.Println(rslt)
		}
	}

	if err := scanner.Err(); err != nil {
		panic("error reading input")
	}
}
