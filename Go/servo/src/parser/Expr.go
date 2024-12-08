package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/helpers"
	"Go/servo/src/lexer"
	"fmt"
	"strconv"
)

func parseExpression(p *parser, bp BindingPower) ast.Expression {
	// Parse NUD
	tokenKind := p.currentTokenKind()
	nudFunction, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s | line:%d:%d ('%s')\n", lexer.TokenKindString(tokenKind), p.getLine(), p.getCharNumber(), p.file))
	}

	left := nudFunction(p)
	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFunction, exists := led_lu[tokenKind]
		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s | line:%d:%d ('%s')\n", lexer.TokenKindString(tokenKind), p.getLine(), p.getCharNumber(), p.file))
		}

		left = ledFunction(p, left, bp_lu[p.currentTokenKind()])

	}
	return left
}

func parsePrimaryExpression(p *parser) ast.Expression {
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

func parseBinaryExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	right := parseExpression(p, bp)

	return ast.BinaryExpr{Left: left, Operator: operatorToken, Right: right}
}

func parsePrefixExpression(p *parser) ast.Expression {
	operatorToken := p.advance()
	rhs := parseExpression(p, default_bp)
	return ast.PrefixExpr{Operator: operatorToken, RightExpr: rhs}
}

func parseGroupingExpression(p *parser) ast.Expression {
	p.advance() // advance past start
	expr := parseExpression(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parseAssignmentExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	rhs := parseExpression(p, bp)
	return ast.AssignmentExpr{Operator: operatorToken, Value: rhs, Assignee: left}
}

func parseStructInstantiationExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	var structName = helpers.ExpectType[ast.SymbolExpr](left).Value
	var properties = map[string]ast.Expression{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		propertyName := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.POINTER)
		expr := parseExpression(p, logical)
		properties[propertyName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.SEMI_COLON)
		}
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.StructInstantiationExpr{StructName: structName, Properties: properties}
}

// TODO: Fix syntax
// Syntax: Number[]
// Current: []Number
func parseArrayInstantiationExpression(p *parser) ast.Expression {
	var underlyingType ast.Type
	var contents []ast.Expression
	var length ast.Expression

	p.expect(lexer.OPEN_BRACKET)
	// Handle sizes defined within the brackets of an array.
	if p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		length = parseExpression(p, logical)
	}

	p.expect(lexer.CLOSE_BRACKET)

	underlyingType = parseType(p, default_bp)

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		contents = append(contents, parseExpression(p, logical)) // Logical prevents assignment within the array
		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.ArrayInstantiationExpr{Underlying: underlyingType, Contents: contents, Length: length}
}
