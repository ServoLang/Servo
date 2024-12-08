package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
)

type BindingPower int

const (
	default_bp BindingPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expression
type led_handler func(p *parser, left ast.Expression, bp BindingPower) ast.Expression

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]BindingPower

// Lookup tables
var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}

func led(kind lexer.TokenKind, bp BindingPower, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	bp_lu[kind] = primary
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {
	led(lexer.ASSIGNMENT, assignment, parseAssignmentExpression)
	led(lexer.PLUS_EQUALS, assignment, parseAssignmentExpression)
	led(lexer.MINUS_EQUALS, assignment, parseAssignmentExpression)
	led(lexer.SLASH_EQUALS, assignment, parseAssignmentExpression)
	led(lexer.STAR_EQUALS, assignment, parseAssignmentExpression)
	led(lexer.MOD_EQUALS, assignment, parseAssignmentExpression)
	led(lexer.POW_EQUALS, assignment, parseAssignmentExpression)

	// Logical Expressions (led expects a value to the left)
	led(lexer.AND, logical, parseBinaryExpression)
	led(lexer.OR, logical, parseBinaryExpression)
	led(lexer.DOT_DOT, logical, parseBinaryExpression)

	// Relational Expressions (led expects a value to the left)
	led(lexer.LESS, relational, parseBinaryExpression)
	led(lexer.LESS_EQUALS, relational, parseBinaryExpression)
	led(lexer.GREATER, relational, parseBinaryExpression)
	led(lexer.GREATER_EQUALS, relational, parseBinaryExpression)
	led(lexer.EQUALS, relational, parseBinaryExpression)
	led(lexer.NOT_EQUALS, relational, parseBinaryExpression)

	// Additive & Multiplicative
	led(lexer.PLUS, additive, parseBinaryExpression)
	led(lexer.DASH, additive, parseBinaryExpression)

	led(lexer.STAR, multiplicative, parseBinaryExpression)
	led(lexer.SLASH, multiplicative, parseBinaryExpression)
	led(lexer.PERCENT, multiplicative, parseBinaryExpression)
	led(lexer.POW, multiplicative, parseBinaryExpression)

	// Literals & Symbols (does not require anything to the left)
	nud(lexer.NUMBER, parsePrimaryExpression)
	nud(lexer.STRING, parsePrimaryExpression)
	nud(lexer.IDENTIFIER, parsePrimaryExpression)
	nud(lexer.OPEN_PAREN, parseGroupingExpression)
	nud(lexer.DASH, parsePrefixExpression)

	// Call / Members / Arrays
	led(lexer.OPEN_CURLY, call, parseStructInstantiationExpression)
	nud(lexer.OPEN_BRACKET, parseArrayInstantiationExpression)

	// Statements
	stmt(lexer.CONST, parseVariableDeclarationStatement)
	stmt(lexer.LET, parseVariableDeclarationStatement)
	stmt(lexer.NEW, parseVariableDeclarationStatement)
	stmt(lexer.STRUCT, parseStructDeclarationStatement)
}
