package scanner

import (
	"fmt"
	"strconv"

	"github.com/ralexgt/glox/errors"
	"github.com/ralexgt/glox/global"
	"github.com/ralexgt/glox/token"
)

type Scanner struct {
	source []rune
	Tokens []*token.Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: []rune(source),
	}
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, &token.Token{
		TokenType: token.TokenType_EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(token.TokenType_LeftParen)
	case ')':
		s.addToken(token.TokenType_RightParen)
	case '{':
		s.addToken(token.TokenType_LeftBrace)
	case '}':
		s.addToken(token.TokenType_RightBrace)
	case ',':
		s.addToken(token.TokenType_Comma)
	case '.':
		s.addToken(token.TokenType_Dot)
	case '-':
		s.addToken(token.TokenType_Minus)
	case '+':
		s.addToken(token.TokenType_Plus)
	case '*':
		s.addToken(token.TokenType_Star)
	case ';':
		s.addToken(token.TokenType_Semicolon)

	case '!':
		if s.match('=') {
			s.addToken(token.TokenType_BangEqual)
		} else {
			s.addToken(token.TokenType_Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.TokenType_EqualEqual)
		} else {
			s.addToken(token.TokenType_Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.TokenType_LessEqual)
		} else {
			s.addToken(token.TokenType_Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.TokenType_GreaterEqual)
		} else {
			s.addToken(token.TokenType_Greater)
		}

	// If the scanner encounters string literals
	case '"':
		s.scanString()

	// If it encounters 2 slashes in a row, consume the rest of the line *it is a comment*
	// /* */ is multiline comment
	case '/':
		if s.match('/') {
			s.consumeLineComment()
			break
		}
		if s.match('*') {
			s.consumeMultiLineComment()
			break
		}
		s.addToken(token.TokenType_Slash)

	case ' ', '\r', '\t':
		break

	case '\n':
		s.line++

	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			global.VM.ReportError(s.line, errors.ErrUnexpectedChar)
		}
	}

}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s *Scanner) advance() rune {
	result := s.source[s.current]
	s.current++
	return result
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		global.VM.ReportError(s.line, errors.ErrUnterminatedString)
		return
	}

	s.advance()

	value := string(s.source[s.start+1 : s.current-1])

	s.addTokenWithLiteral(token.TokenType_String, value)
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		// consume the .
		s.advance()

		// consume the digits after the floating point
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	text := string(s.source[s.start:s.current])
	num, err := strconv.ParseFloat(text, 64)
	if err != nil {
		global.VM.ReportError(s.line, errors.ErrInvalidNumber)
		return
	}

	s.addTokenWithLiteral(token.TokenType_Number, num)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	if t, ok := token.Keywords[text]; ok {
		s.addToken(t)
		return
	}

	s.addToken(token.TokenType_Identifier)
}

func (s *Scanner) addToken(t token.TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t token.TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.Tokens = append(s.Tokens, &token.Token{
		TokenType: t,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.line,
	})
}

func (s *Scanner) consumeLineComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) consumeMultiLineComment() {
	for !s.isAtEnd() {
		switch c := s.peek(); c {
		case '\n':
			s.line++
			s.advance()

		case '*':
			if s.peekNext() == '/' {
				s.current += 2
				return
			}
			s.advance()

		default:
			s.advance()
		}
	}

	global.VM.ReportError(s.line, errors.ErrUnterminatedComment)
}

// Moved from global to solve circular dependency
func (l *global.Lox) Run(source string) {
	scanner := NewScanner(source)
	scanner.ScanTokens()

	for _, token := range scanner.Tokens {
		fmt.Println(token)
	}
}
