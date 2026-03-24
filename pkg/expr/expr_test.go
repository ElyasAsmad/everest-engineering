package expr_test

import (
	"testing"

	"github.com/ElyasAsmad/everestengineering2/pkg/expr"
)

func TestLexer_BasicTokens(t *testing.T) {
	input := "d < 200 && 70 <= w <= 200"
	lexer := expr.NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("unexpected lexer error: %v", err)
	}

	types := make([]expr.TokenType, len(tokens))
	for i, token := range tokens {
		types[i] = token.Type
	}

	expectedTypes := []expr.TokenType{
		expr.TOKEN_IDENT,  // d
		expr.TOKEN_LT,     // <
		expr.TOKEN_NUMBER, // 200
		expr.TOKEN_AND,    // &&
		expr.TOKEN_NUMBER, // 70
		expr.TOKEN_LTE,    // <=
		expr.TOKEN_IDENT,  // w
		expr.TOKEN_LTE,    // <=
		expr.TOKEN_NUMBER, // 200
		expr.TOKEN_EOF,    // EOF
	}

	if len(types) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d\n got: %v want: %v", len(expectedTypes), len(types), types, expectedTypes)
	}

	for i := range expectedTypes {
		if types[i] != expectedTypes[i] {
			t.Errorf("token[%d]: got %s, want %s", i, types[i], expectedTypes[i])
		}
	}
}

func TestLexer_AndKeyword(t *testing.T) {
	lexer := expr.NewLexer("d < 200 and w > 50")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	// try to find the AND token
	found := false
	for _, token := range tokens {
		if token.Type == expr.TOKEN_AND {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected TOKEN_AND for keyword 'and'")
	}
}

func TestLexer_InvalidChar(t *testing.T) {
	lexer := expr.NewLexer("d @ 200")
	_, err := lexer.Tokenize()
	if err == nil {
		t.Fatal("expected error for '@', got nil")
	}
}

func TestLexer_SingleEquals(t *testing.T) {
	lexer := expr.NewLexer("d = 200")
	_, err := lexer.Tokenize()
	if err == nil {
		t.Fatal("expected error for '=', got nil")
	}
}

// --- Parser + Evaluator tests ---

type evalCase struct {
	name string
	expr string
	ctx  expr.Context
	want bool
}

func runEvalCases(t *testing.T, cases []evalCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := expr.Compile(tc.expr)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}
			got, err := expr.Evaluate(node, tc.ctx)
			if err != nil {
				t.Fatalf("evaluate error: %v", err)
			}
			if got != tc.want {
				t.Errorf("expr %q with ctx %v: got %v, want %v", tc.expr, tc.ctx, got, tc.want)
			}
		})
	}
}

func TestSimpleComparisons(t *testing.T) {
	runEvalCases(t, []evalCase{
		{"lt_true", "d < 200", expr.Context{"d": 150}, true},
		{"lt_false", "d < 200", expr.Context{"d": 200}, false},
		{"lte_true", "d <= 200", expr.Context{"d": 200}, true},
		{"gt_true", "w > 50", expr.Context{"w": 100}, true},
		{"gt_false", "w > 50", expr.Context{"w": 50}, false},
		{"gte_true", "w >= 50", expr.Context{"w": 50}, true},
		{"eq_true", "d == 100", expr.Context{"d": 100}, true},
		{"eq_false", "d == 100", expr.Context{"d": 101}, false},
	})
}

func TestChainedComparisons(t *testing.T) {
	runEvalCases(t, []evalCase{
		{"in_range", "70 <= w <= 200", expr.Context{"w": 100}, true},
		{"on_lower_bound", "70 <= w <= 200", expr.Context{"w": 70}, true},
		{"on_upper_bound", "70 <= w <= 200", expr.Context{"w": 200}, true},
		{"below_range", "70 <= w <= 200", expr.Context{"w": 69}, false},
		{"above_range", "70 <= w <= 200", expr.Context{"w": 201}, false},
		{"strictly_in_range", "70 < w < 200", expr.Context{"w": 70}, false},
	})
}

func TestAndExpressions(t *testing.T) {
	runEvalCases(t, []evalCase{
		{"ofr001_ok", "d < 200 && 70 <= w <= 200", expr.Context{"d": 150, "w": 100}, true},
		{"ofr001_fail_distance", "d < 200 && 70 <= w <= 200", expr.Context{"d": 250, "w": 100}, false},
		{"ofr001_fail_weight", "d < 200 && 70 <= w <= 200", expr.Context{"d": 150, "w": 50}, false},
		{"ofr001_fail_both", "d < 200 && 70 <= w <= 200", expr.Context{"d": 250, "w": 50}, false},
		// using 'and' instead of '&&'
		{"ofr001_and_ok", "d < 200 and 70 <= w <= 200", expr.Context{"d": 150, "w": 100}, true},
		{"ofr001_and_fail_distance", "d < 200 and 70 <= w <= 200", expr.Context{"d": 250, "w": 100}, false},
		{"ofr001_and_fail_weight", "d < 200 and 70 <= w <= 200", expr.Context{"d": 150, "w": 50}, false},
		{"ofr001_and_fail_both", "d < 200 and 70 <= w <= 200", expr.Context{"d": 250, "w": 50}, false},
	})
}

func TestMultipleAndChain(t *testing.T) {
	runEvalCases(t, []evalCase{
		{"three_clauses_ok", "d < 300 && w > 10 && w < 200", expr.Context{"d": 100, "w": 100}, true},
		{"three_clauses_fail", "d < 300 && w > 10 && w < 200", expr.Context{"d": 100, "w": 5}, false},
	})
}

func TestUndefinedVariable(t *testing.T) {
	node, err := expr.Compile("x < 100")
	if err != nil {
		t.Fatal(err)
	}
	_, err = expr.Evaluate(node, expr.Context{"d": 50}) // 'x' is not defined in context
	if err == nil {
		t.Fatal("expected error for undefined variable 'x', got nil")
	}
}

func TestParseErrors(t *testing.T) {
	bad := []string{
		"< 200",      // missing left operand
		"d < ",       // missing right operand
		"d < 200 &&", // missing right operand for second clause
	}
	for _, expr_str := range bad {
		_, err := expr.Compile(expr_str)
		if err == nil {
			t.Fatalf("expected error for expression %q, got nil", expr_str)
		}
	}
}
