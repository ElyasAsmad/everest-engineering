package model

import "github.com/ElyasAsmad/everestengineering2/pkg/expr"

type Offer struct {
	Code       string
	Discount   float64
	Constraint string
	compiled   expr.Node
}
