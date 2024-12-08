package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
	"fmt"
	"strconv"
)

func parseExpression(p *parser, bp binding_power) ast.Expr {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}

func parsePrefixExpression(p *parser) ast.Expr {
	operatorToken := p.advance()
	expr := parseExpression(p, unary)

	return ast.PrefixExpr{
		Operator: operatorToken,
		Right:    expr,
	}
}

func parseAssignmentExpression(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()
	rhs := parseExpression(p, bp)

	return ast.AssignmentExpr{
		Assigne:       left,
		AssignedValue: rhs,
	}
}

func parseRangeExpression(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()

	return ast.RangeExpr{
		Lower: left,
		Upper: parseExpression(p, bp),
	}
}

func parseBinaryExpression(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	right := parseExpression(p, defalt_bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

func parsePrimaryExpression(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpr{
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: p.advance().Value,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expr from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseMemberExpression(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	isComputed := p.advance().Kind == lexer.OPEN_BRACKET

	if isComputed {
		rhs := parseExpression(p, bp)
		p.expect(lexer.CLOSE_BRACKET)
		return ast.ComputedExpr{
			Member:   left,
			Property: rhs,
		}
	}

	return ast.MemberExpr{
		Member:   left,
		Property: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parseArrayLiteralExpression(p *parser) ast.Expr {
	p.expect(lexer.OPEN_BRACKET)
	arrayContents := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		arrayContents = append(arrayContents, parseExpression(p, logical))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_BRACKET) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_BRACKET)

	return ast.ArrayLiteral{
		Contents: arrayContents,
	}
}

func parseGroupingExpression(p *parser) ast.Expr {
	p.expect(lexer.OPEN_PAREN)
	expr := parseExpression(p, defalt_bp)
	p.expect(lexer.OPEN_PAREN)
	return expr
}

func parseCallExpression(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()
	arguments := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		arguments = append(arguments, parseExpression(p, assignment))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_PAREN) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)
	return ast.CallExpr{
		Method:    left,
		Arguments: arguments,
	}
}

func parseFunctionExpression(p *parser) ast.Expr {
	p.expect(lexer.FUNCTION)
	functionParams, returnType, functionBody := parseFunctionParamsAndBody(p)

	return ast.FunctionExpr{
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
	}
}
