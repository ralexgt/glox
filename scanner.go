package main

import "strconv"

type Scanner struct {
	source []rune
	tokens []*Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: []rune(source),
	}
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, &Token{
		TokenType: TokenType_EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(TokenType_LeftParen)
	case ')':
		s.addToken(TokenType_RightParen)
	case '{':
		s.addToken(TokenType_LeftBrace)
	case '}':
		s.addToken(TokenType_RightBrace)
	case ',':
		s.addToken(TokenType_Comma)
	case '.':
		s.addToken(TokenType_Dot)
	case '-':
		s.addToken(TokenType_Minus)
	case '+':
		s.addToken(TokenType_Plus)
	case '*':
		s.addToken(TokenType_Star)
	case ';':
		s.addToken(TokenType_Semicolon)

	case '!':
		if s.match('=') {
			s.addToken(TokenType_BangEqual)
		} else {
			s.addToken(TokenType_Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(TokenType_EqualEqual)
		} else {
			s.addToken(TokenType_Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(TokenType_LessEqual)
		} else {
			s.addToken(TokenType_Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(TokenType_GreaterEqual)
		} else {
			s.addToken(TokenType_Greater)
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
		s.addToken(TokenType_Slash)

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
			vm.reportError(s.line, ErrUnexpectedChar)
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
		vm.reportError(s.line, ErrUnterminatedString)
		return
	}

	s.advance()

	value := string(s.source[s.start+1 : s.current-1])

	s.addTokenWithLiteral(TokenType_String, value)
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
		vm.reportError(s.line, ErrInvalidNumber)
		return
	}

	s.addTokenWithLiteral(TokenType_Number, num)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	if t, ok := keywords[text]; ok {
		s.addToken(t)
		return
	}

	s.addToken(TokenType_Identifier)
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, &Token{
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

	vm.reportError(s.line, ErrUnterminatedComment)
}
