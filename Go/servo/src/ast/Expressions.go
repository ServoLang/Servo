package ast

import (
	"Go/servo/src/lexer"
)

// ===================================================================================================================\\
// LITERAL EXPRESSIONS
// ===================================================================================================================\\

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

// ===================================================================================================================\\
// END LITERAL EXPRESSIONS
// ===================================================================================================================\\

// ===================================================================================================================\\
// COMPLEX EXPRESSIONS
// ===================================================================================================================\\

type BinaryExpr struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

func (n BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expression
}

func (n PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assignee Expression
	Operator lexer.Token
	Value    Expression
}

func (n AssignmentExpr) expr() {}

type StructInstantiationExpr struct {
	StructName string
	Properties map[string]Expression
}

func (n StructInstantiationExpr) expr() {}

type ArrayInstantiationExpr struct {
	Underlying Type
	Length     Expression
	Contents   []Expression
}

func (n ArrayInstantiationExpr) expr() {}

// ===================================================================================================================\\
// END COMPLEX EXPRESSIONS
// ===================================================================================================================\\
