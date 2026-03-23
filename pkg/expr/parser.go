package expr

import (
	"fmt"
	"strconv"
)

// converts a flat token list into AST
type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (Node, error) {
	node, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.isEOF() {
		tok := p.current()
		return nil, fmt.Errorf("pos %d: unexpected token %q after end of expression", tok.Pos, tok.Literal)
	}
	return node, nil
}

// grammar (precedence from lowest to highest):
// expr -> comparison (AND comparison)*
// comparison -> operand (compOp operand (compOp operand)?)?
// operand -> NUMBER | IDENT
// compOp -> < | <= | > | >= | ==

func (p *Parser) parseExpr() (Node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.currentType() == TOKEN_AND {
		p.advance() // consume '&&'

		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		left = &AndNode{
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) parseComparison() (Node, error) {
	left, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	if !p.isCompOp(p.currentType()) {
		return left, nil // no comparison operator, just return the operand
	}

	op1, err := tokenTypeToOperator(p.currentType())
	if err != nil {
		return nil, err
	}
	p.advance()

	mid, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	// check for chained comparison: left op1 mid op2 right
	if p.isCompOp(p.currentType()) {
		op2, err := tokenTypeToOperator(p.currentType())
		if err != nil {
			return nil, err
		}
		p.advance() // consume second operator

		right, err := p.parseOperand()
		if err != nil {
			return nil, err
		}

		return &ChainedCompareNode{
			Left:  left,
			Op1:   op1,
			Mid:   mid,
			Op2:   op2,
			Right: right,
		}, nil
	}

	return &CompareNode{
		Left:  left,
		Op:    op1,
		Right: mid,
	}, nil
}

func (p *Parser) parseOperand() (Node, error) {
	tok := p.current()

	switch tok.Type {
	case TOKEN_NUMBER:
		val, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("pos %d: invalid number: %q", tok.Pos, tok.Literal)
		}
		p.advance()
		return &NumberNode{Value: val}, nil
	case TOKEN_IDENT:
		p.advance()
		return &VariableNode{Name: tok.Literal}, nil

	case TOKEN_EOF:
		return nil, fmt.Errorf("pos %d: unexpected end of expression, expected operand", tok.Pos)

	default:
		return nil, fmt.Errorf("pos %d: unexpected number or variable, got %q", tok.Pos, tok.Literal)
	}
}

// current() returns the current token, or EOF token if we have reached the end of input
func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) currentType() TokenType {
	return p.current().Type
}

// advance() moves to the next token. If we are already at the end, it does nothing.
func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) isEOF() bool {
	return p.currentType() == TOKEN_EOF
}

func (p *Parser) isCompOp(t TokenType) bool {
	switch t {
	case TOKEN_LT, TOKEN_LTE, TOKEN_GT, TOKEN_GTE, TOKEN_EQ:
		return true
	}
	return false
}
