package scan

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/brentellingson/go-lox/internal/token"
)

var reserved = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"fun":    token.FUN,
	"for":    token.FOR,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type ScanError struct {
	Line    int
	Message string
}

func NewScanError(line int, message string) *ScanError {
	return &ScanError{Line: line, Message: message}
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("Scan Error on Line %v: %v", e.Line, e.Message)
}

func Scan(source string) ([]token.Token, error) {
	scanner := NewScanner(source)
	return scanner.ScanTokens(), nil
}

type Scanner struct {
	Source  string
	Tokens  []token.Token
	start   int
	current int
	line    int
	errs    []error
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source: source,
		line:   1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.Tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!': // BANG or BANG_EQUAL
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=': // EQUAL or EQUAL_EQUAL
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<': // LESS or LESS_EQUAL
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>': // GREATER or GREATER_EQUAL
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '"':
		s.string()
	case '\n':
		s.line++
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.errs = append(s.errs, NewScanError(s.line, "Unexpected character "+string(c)))
		}
	}
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errs = append(s.errs, NewScanError(s.line, "Unterminated string."))
		return
	}

	s.advance() // the closing ".

	value := s.Source[s.start+1 : s.current-1]
	s.addTokenLiteral(token.STRING, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance() // consume the "."
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.Source[s.start:s.current], 64)
	if err != nil {
		s.errs = append(s.errs, NewScanError(s.line, "unable to parse number "+s.Source[s.start:s.current]))
		return
	}
	s.addTokenLiteral(token.NUMBER, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	value := s.Source[s.start:s.current]
	if keyword, ok := reserved[value]; ok {
		s.addToken(keyword)
		return
	}

	s.addToken(token.IDENTIFIER)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) advance() rune {
	r, size := utf8.DecodeRuneInString(s.Source[s.current:])
	s.current += size
	return r
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	r, size := utf8.DecodeRuneInString(s.Source[s.current:])
	if r != expected {
		return false
	}
	s.current += size
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(s.Source[s.current:])
	return r
}

func (s *Scanner) peekNext() rune {
	if s.isAtEnd() {
		return 0
	}
	_, size := utf8.DecodeRuneInString(s.Source[s.current:])
	if s.current+size >= len(s.Source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s.Source[s.current+size:])
	return r
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, token.NewToken(tokenType, text, nil, s.line))
}

func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal any) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, token.NewToken(tokenType, text, literal, s.line))
}
