package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
	"fmt"
	"strconv"
)

func parse_expr(p *parser, bp binding_power) ast.Expr {
	// Parse NUD
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf(`NUD HANDLER EXPECTED FOR TOKEN %s\n`, lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)
	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]
		if !exists {
			panic(fmt.Sprintf(`LED HANDLER EXPECTED FOR TOKEN %s\n`, lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp_lu[p.currentTokenKind()])

	}
	return left
}

func parse_primary_expr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{Value: number}
	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{Value: p.advance().Value}
	default:
		panic(fmt.Sprintf(`Cannot create primary_expression from %s\n`, lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parse_binary_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	right := parse_expr(p, bp)

	return ast.BinaryExpr{Left: left, Operator: operatorToken, Right: right}
}

func parse_prefix_expr(p *parser) ast.Expr {
	operatorToken := p.advance()
	rhs := parse_expr(p, default_bp)
	return ast.PrefixExpr{Operator: operatorToken, RightExpr: rhs}
}

func parse_grouping_expr(p *parser) ast.Expr {
	p.advance() // advance past start
	expr := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parse_assignment_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	rhs := parse_expr(p, bp)
	return ast.AssignmentExpr{Operator: operatorToken, Value: rhs, Assignee: left}
}
