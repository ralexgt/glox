package scanner

import (
	"strconv"

	"github.com/ralexgt/glox/errors"
	"github.com/ralexgt/glox/token"
)

// Scanner will use this to report errors
type ErrorHandler func(line int, err error)

type Scanner struct {
	source []rune
	Tokens []*token.Token

	start   int
	current int
	line    int

	errorHandler ErrorHandler // Function to handle errors
}

// NewScanner takes an ErrorHandler function as an argument
func NewScanner(source string, errorHandler ErrorHandler) *Scanner {
	return &Scanner{
		source:       []rune(source),
		line:         0,            // Initialize line counter
		errorHandler: errorHandler, // Store the provided error handler
	}
}

func (s *Scanner) ScanTokens() {
	s.Tokens = nil // Clear any previous tokens
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			// Error already handled by errorHandler within scanToken, stop scanning
			return
		}
	}

	s.Tokens = append(s.Tokens, &token.Token{
		TokenType: token.TokenType_EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
}

func (s *Scanner) scanToken() error {
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

	case '"':
		if err := s.scanString(); err != nil {
			return err
		}

	case '/':
		if s.match('/') {
			s.consumeLineComment()
			break
		}
		if s.match('*') {
			if err := s.consumeMultiLineComment(); err != nil {
				return err
			}
			break
		}
		s.addToken(token.TokenType_Slash)

	case ' ', '\r', '\t':
		// Ignore whitespace
		break

	case '\n':
		s.line++ // Increment line number on newline

	default:
		if isDigit(c) {
			if err := s.scanNumber(); err != nil {
				return err
			}
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			s.errorHandler(s.line, errors.ErrUnexpectedChar) // Use injected error handler
		}
	}

	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0 // Null rune if at end
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
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) scanString() error {
	for ; s.peek() != '"' && !s.isAtEnd(); s.advance() {
		if s.peek() == '\n' {
			s.line++ // Allow newlines in strings
		}
	}

	if s.isAtEnd() {
		s.errorHandler(s.line, errors.ErrUnterminatedString) // Use injected error handler
		return errors.ErrUnterminatedString
	}

	s.advance() // Consume the closing "

	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(token.TokenType_String, value)
	return nil
}

func (s *Scanner) scanNumber() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance() // Consume the "."
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	numberStr := string(s.source[s.start:s.current])
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		s.errorHandler(s.line, errors.ErrInvalidNumber) // Use injected error handler
		return errors.ErrInvalidNumber
	}

	s.addTokenWithLiteral(token.TokenType_Number, number)
	return nil
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	if tokType, isKeyword := token.Keywords[text]; isKeyword {
		s.addToken(tokType)
	} else {
		s.addToken(token.TokenType_Identifier)
	}
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.Tokens = append(s.Tokens, &token.Token{
		TokenType: tokenType,
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

func (s *Scanner) consumeMultiLineComment() error {
	for !s.isAtEnd() {
		if s.peek() == '*' && s.peekNext() == '/' {
			s.current += 2 // Consume "*/"
			return nil     // Successfully consumed multiline comment
		}
		if s.peek() == '\n' {
			s.line++ // Increment line number for newlines in comments
		}
		s.advance() // Consume character within comment
	}

	s.errorHandler(s.line, errors.ErrUnterminatedComment) // Use injected error handler
	return errors.ErrUnterminatedComment
}
