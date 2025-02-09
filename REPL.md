# Read-Eval-Print Loop for Lox

```
for {
    tokens := []
    statements := []
    for {
        source := read()
        tokens.append scan(source)
        for {
             statements.append parse(tokens)
        }
        if 
    }
    result := execute(statements)
    print(result)
}

[text](internal/repl/repl2.go)


[text](internal/repl/repl2.go)
// package repl implements the Read-Eval-Print-Loop for the Lox Programming Language.
package repl

import (
	"errors"

	"github.com/brentellingson/go-lox/internal/token"
)

var (
	ParseError          = errors.New("parse error")
	IncompleteStatement = errors.New("incomplete statement")
)

type TokenBuffer interface {
	Append(tokens []token.Token)
	IsAtEnd() bool
}

// type Reader interface {
// 	Read(prompt string) ([]rune, error)
// }

// type Scanner interface {
// 	Scan(text []rune) ([]token.Token, error)
// }

// type Parser interface {
// 	Parse(tokens TokenBuffer) (ast.Stmt, error)
// 	Synchronize(tokens TokenBuffer)
// }

// type Interpreter interface {
// 	Interpret(statements []internal.Stmt) any
// }

// Read-Eval-Print does this:
// while !bailout && !done{}
//    text, err = read()
//    tokens, err = scan(text)
//    if
//    tokenbuff.append(tokens)
// }
// Read One Line
// Scan it into tokens
// if can't scan, print error and return
// Parse into a statement
// if not bailoout and incomplete statment, read and tokenize another line
// if syntax error, read to next sync point

// type Repl struct {
// 	reader      Reader
// 	scanner     Scanner
// 	buffer      func([]token.Token) TokenBuffer
// 	parser      Parser
// 	interpreter Interpreter
// }

// func New() *Repl {
// 	return &Repl{}
// }

// // Read reads a line of input, scans it into tokens, parses the tokens into statements, and returns the statements.
// // If there are any errors, it returns the errors.
// // If the input is incomplete, it reads more input and continues parsing.
// func (r *Repl) Read() ([]internal.Stmt, error) {
// 	var errs []error
// 	var stmts []internal.Stmt

// 	text, err := r.reader.Read(">>> ")
// 	if err != nil {
// 		return nil, err
// 	}
// 	tokens, err := r.scanner.Scan(text)
// 	if err != nil {
// 		return nil, err
// 	}
// 	buffer := r.buffer(tokens)

// 	for !buffer.IsAtEnd() {
// 		stmt, err := r.parser.Parse(buffer)
// 		if errors.Is(err, IncompleteStatement) && len(errs) == 0 {
// 			// incomplete statement but no syntax errors yet; ask for more input and keep going
// 			text, err2 := r.reader.Read("... ")
// 			if err2 == io.EOF {
// 				// user hit Ctrl-D/Ctrl-Z; stop parsing
// 				errs = append(errs, err)
// 				break
// 			}
// 			if err2 != nil {
// 				// some other error; stop parsing
// 				errs = append(errs, err, err2)
// 				break
// 			}
// 			tokens, err2 := r.scanner.Scan(text)
// 			if err2 != nil {
// 				// couldn't scan the input; stop parsing
// 				errs = append(errs, err2)
// 				break
// 			}
// 			buffer.Append(tokens)
// 			continue
// 		}
// 		if errors.Is(err, ParseError) {
// 			// syntax error; synchronize and try keep going
// 			errs = append(errs, err)
// 			r.parser.Synchronize(buffer)
// 			continue
// 		}
// 		if err != nil {
// 			// some other error; stop parsing
// 			errs = append(errs, err)
// 			break
// 		}
// 		if stmt != nil {
// 			stmts = append(stmts, stmt)
// 		}
// 	}
// 	if len(errs) > 0 {
// 		return nil, errors.Join(errs...)
// 	}
// 	return stmts, nil
// }
