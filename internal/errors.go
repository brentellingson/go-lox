package internal

import "fmt"

var (
	HadParserError  = false
	HadRuntimeError = false
)

func ReportParserError(line int, message string) {
	fmt.Printf("[line %d] Parse Error: %s\n", line, message)
	HadParserError = true
}

func ReportRuntimeError(err error) {
	fmt.Println(err.Error())
	HadRuntimeError = true
}
