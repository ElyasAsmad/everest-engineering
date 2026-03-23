package expr

import (
	"fmt"
	"strconv"
	"strings"
)

// base interface for all AST nodes
type Node interface {
	node()
	String() string
}

// leaf nodes
type NumberNode struct {
	Value float64
}

func (n *NumberNode) node() {}
func (n *NumberNode) String() string {
	return strconv.FormatFloat(n.Value, 'f', -1, 64)
}

type VariableNode struct {
	Name string
}

func (v *VariableNode) node() {}
func (v *VariableNode) String() string {
	return v.Name
}

// operator nodes
type Operator string

const (
	OpLT  Operator = "<"
	OpLTE Operator = "<="
	OpGT  Operator = ">"
	OpGTE Operator = ">="
	OpEQ  Operator = "=="
)

func tokenTypeToOperator(t TokenType) (Operator, error) {
	switch t {
	case TOKEN_LT:
		return OpLT, nil
	case TOKEN_LTE:
		return OpLTE, nil
	case TOKEN_GT:
		return OpGT, nil
	case TOKEN_GTE:
		return OpGTE, nil
	case TOKEN_EQ:
		return OpEQ, nil
	default:
		return "", fmt.Errorf("token %s is not a comparison operator", t)
	}
}

// comparison node

// handles simple binary comparisons like "weight > 10" or "distance <= 5"
type CompareNode struct {
	Left  Node
	Op    Operator
	Right Node
}

func (c *CompareNode) node() {}
func (c *CompareNode) String() string {
	return fmt.Sprintf("(%s %s %s)", c.Left, c.Op, c.Right)
}

// ChainedCompareNode handles chained comparisons like "10 < weight <= 20"
type ChainedCompareNode struct {
	Left  Node     // e.g. 10
	Op1   Operator // e.g. <
	Mid   Node     // e.g. weight
	Op2   Operator // e.g. <=
	Right Node     // e.g. 20
}

func (c *ChainedCompareNode) node() {}
func (c *ChainedCompareNode) String() string {
	return fmt.Sprintf("(%s %s %s %s %s)", c.Left, c.Op1, c.Mid, c.Op2, c.Right)
}

// Logical nodes
// handles logical AND of two comparisons, e.g. "(weight > 10) AND (distance <= 5)"
type AndNode struct {
	Left  Node
	Right Node
}

func (a *AndNode) node() {}
func (a *AndNode) String() string {
	return fmt.Sprintf("(%s && %s)", a.Left, a.Right)
}

// helpers
func Dump(n Node, indent int) string {
	prefix := strings.Repeat(" ", indent)
	switch v := n.(type) {
	case *NumberNode:
		return fmt.Sprintf("%sNumber(%s)\n", prefix, v)
	case *VariableNode:
		return fmt.Sprintf("%sVariable(%s)\n", prefix, v.Name)
	case *CompareNode:
		return fmt.Sprintf("%sCompare(%s)\n%s%s", prefix, v.Op, Dump(v.Left, indent+1), Dump(v.Right, indent+1))
	case *ChainedCompareNode:
		return fmt.Sprintf("%sChainedCompare(%s, %s)\n%s%s%s", prefix, v.Op1, v.Op2, Dump(v.Left, indent+1), Dump(v.Mid, indent+1), Dump(v.Right, indent+1))
	case *AndNode:
		return fmt.Sprintf("%sAnd\n%s%s", prefix, Dump(v.Left, indent+1), Dump(v.Right, indent+1))
	default:
		return fmt.Sprintf("%sUnknown(%T)\n", prefix, n)
	}
}
