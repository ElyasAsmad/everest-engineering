package expr

import "strings"

// Compile is a convenience function: lexes + parses an expression string
// returns a ready-to-evaluate AST node

// example usage:
// node, err := expr.Compile("d < 200 && 70 <= w <= 200")
// ok, err := expr.Evaluate(node, expr.Context{"d": 150, "w": 80})
func Compile(expression string) (Node, error) {
	// if the expression is empty or just whitespace (no constraints) = always true
	if strings.TrimSpace(expression) == "" {
		return &BoolNode{Value: true}, nil
	}

	lexer := NewLexer(expression)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	return parser.Parse()
}
