package expr

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	// Literals & Identifiers
	TOKEN_NUMBER TokenType = iota
	TOKEN_IDENT

	// comparison operators
	TOKEN_LT  // <
	TOKEN_LTE // <=
	TOKEN_GT  // >
	TOKEN_GTE // >=
	TOKEN_EQ  // ==

	// Logical operators
	TOKEN_AND // &&

	// Control
	TOKEN_EOF
	TOKEN_ILLEGAL
)

func (t TokenType) String() string {
	switch t {
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_IDENT:
		return "IDENT"
	case TOKEN_LT:
		return "<"
	case TOKEN_LTE:
		return "<="
	case TOKEN_GT:
		return ">"
	case TOKEN_GTE:
		return ">="
	case TOKEN_EQ:
		return "=="
	case TOKEN_AND:
		return "&&"
	case TOKEN_EOF:
		return "EOF"
	default:
		return "ILLEGAL"
	}
}

// single lexical unit
type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %q, pos=%d)", t.Type, t.Literal, t.Pos)
}

// lexer for tokenizing input string
type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
	}
}

// peek() returns the current character without advancing the position.
// It returns false if we have reached the end of input.
func (l *Lexer) peek() (rune, bool) {
	// if pos is out of bounds, return 0 and false
	if l.pos >= len(l.input) {
		return 0, false
	}
	return l.input[l.pos], true
}

func (l *Lexer) peekAt(offset int) (rune, bool) {
	i := l.pos + offset
	if i >= len(l.input) {
		return 0, false
	}
	return l.input[i], true
}

// advance() returns the current character and moves the position forward by one.
func (l *Lexer) advance() rune {
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) skipWhitespace() {
	for {
		ch, ok := l.peek()
		if !ok || !unicode.IsSpace(ch) {
			break
		}
		l.advance()
	}
}

// Tokenize converts the input string into a slice of tokens. It returns an error if it encounters an illegal character or an invalid token.
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}
	return tokens, nil
}

func (l *Lexer) nextToken() (Token, error) {
	l.skipWhitespace()

	startPos := l.pos
	ch, ok := l.peek()
	// if we have reached end of input, return EOF token
	if !ok {
		return Token{Type: TOKEN_EOF, Pos: startPos}, nil
	}

	// handle numbers
	if unicode.IsDigit(ch) {
		return l.readNumber(startPos), nil
	}

	// handle identifiers
	if unicode.IsLetter(ch) {
		return l.readIdentifier(startPos), nil
	}

	l.advance()

	switch ch {
	case '<':
		if next, ok := l.peekAt(0); ok && next == '=' {
			l.advance()
			return Token{Type: TOKEN_LTE, Literal: "<=", Pos: startPos}, nil
		}
		return Token{Type: TOKEN_LT, Literal: "<", Pos: startPos}, nil
	case '>':
		if next, ok := l.peekAt(0); ok && next == '=' {
			l.advance()
			return Token{Type: TOKEN_GTE, Literal: ">=", Pos: startPos}, nil
		}
		return Token{Type: TOKEN_GT, Literal: ">", Pos: startPos}, nil
	case '=':
		if next, ok := l.peekAt(0); ok && next == '=' {
			l.advance()
			return Token{Type: TOKEN_EQ, Literal: "==", Pos: startPos}, nil
		}
		return Token{Type: TOKEN_ILLEGAL, Literal: string(ch), Pos: startPos}, fmt.Errorf("pos %d: single '=' is not a valid operator, did you mean '=='?", startPos)
	case '&':
		if next, ok := l.peekAt(0); ok && next == '&' {
			l.advance()
			return Token{Type: TOKEN_AND, Literal: "&&", Pos: startPos}, nil
		}
		return Token{Type: TOKEN_ILLEGAL, Literal: string(ch), Pos: startPos}, fmt.Errorf("pos %d: single '&' is not a valid operator, did you mean '&&'?", startPos)
	}

	return Token{Type: TOKEN_ILLEGAL, Literal: string(ch), Pos: startPos}, fmt.Errorf("pos %d: unexpected character '%q'", startPos, ch)
}

func (l *Lexer) readNumber(startPos int) Token {
	var sb strings.Builder
	for {
		ch, ok := l.peek()
		// if reached end of input or current character is not a digit or dot, break
		if !ok || (!unicode.IsDigit(ch) && ch != '.') {
			break
		}
		sb.WriteRune(l.advance())
	}
	return Token{Type: TOKEN_NUMBER, Literal: sb.String(), Pos: startPos}
}

func (l *Lexer) readIdentifier(startPos int) Token {
	var sb strings.Builder
	for {
		ch, ok := l.peek()
		// if reached end of input or current character is not a letter, digit or underscore, break
		if !ok || (!unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_') {
			break
		}
		sb.WriteRune(l.advance())
	}
	literal := sb.String()

	// `and` is a keyword alias for `&&`
	if strings.ToLower(literal) == "and" {
		return Token{Type: TOKEN_AND, Literal: "&&", Pos: startPos}
	}

	return Token{Type: TOKEN_IDENT, Literal: literal, Pos: startPos}
}
