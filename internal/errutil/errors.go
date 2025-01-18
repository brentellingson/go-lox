package errutil

import "fmt"

func Error(line int, message string) {
	Report(line, "", message)
}

var HadError = false

func Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	HadError = true
}
