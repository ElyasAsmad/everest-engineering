package expr

import "fmt"

// holds the variable values for evaluation (e.g.: {"d": 150, "w": 80})
type Context map[string]float64

// evaluates the AST with given variable values
func Evaluate(node Node, ctx Context) (bool, error) {
	switch n := node.(type) {
	case *AndNode:
		left, err := Evaluate(n.Left, ctx)
		if err != nil {
			return false, err
		}

		// short-circuit evaluation: if left is false, no need to evaluate right
		if !left {
			return false, nil
		}

		return Evaluate(n.Right, ctx)
	case *CompareNode:
		left, err := resolveValue(n.Left, ctx)
		if err != nil {
			return false, err
		}
		right, err := resolveValue(n.Right, ctx)
		if err != nil {
			return false, err
		}
		return applyOp(left, n.Op, right)
	case *ChainedCompareNode:
		left, err := resolveValue(n.Left, ctx)
		if err != nil {
			return false, err
		}
		mid, err := resolveValue(n.Mid, ctx)
		if err != nil {
			return false, err
		}
		right, err := resolveValue(n.Right, ctx)
		if err != nil {
			return false, err
		}

		firstHalf, err := applyOp(left, n.Op1, mid)
		if err != nil {
			return false, nil
		}
		if !firstHalf {
			return false, nil
		}
		return applyOp(mid, n.Op2, right)
	case *BoolNode:
		return n.Value, nil
	default:
		return false, fmt.Errorf("evaluate: unexpected node type %T", node)
	}
}

// resolveValue extracts a float64 from either NumberNode / VariableNode
func resolveValue(node Node, ctx Context) (float64, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, nil
	case *VariableNode:
		val, ok := ctx[n.Name]
		if !ok {
			return 0, fmt.Errorf("undefined variable %q", n.Name)
		}
		return val, nil
	default:
		return 0, fmt.Errorf("resolveValue: expected leaf node, got %T", node)
	}
}

// applyOp applies the comparison operator to two float64 values
func applyOp(left float64, op Operator, right float64) (bool, error) {
	switch op {
	case OpLT:
		return left < right, nil
	case OpLTE:
		return left <= right, nil
	case OpGT:
		return left > right, nil
	case OpGTE:
		return left >= right, nil
	case OpEQ:
		return left == right, nil
	default:
		return false, fmt.Errorf("unknown operator: %v", op)
	}
}
