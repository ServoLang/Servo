package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
)

type binding_power int

const (
	defalt_bp binding_power = iota
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
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp binding_power) ast.Expr

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]binding_power

var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}

func led(kind lexer.TokenKind, bp binding_power, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, bp binding_power, nud_fn nud_handler) {
	bp_lu[kind] = primary
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	bp_lu[kind] = defalt_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {
	// Assignment
	led(lexer.ASSIGNMENT, assignment, parseAssignmentExpression)
	led(lexer.PLUS_EQUALS, assignment, parseAssignmentExpression)
	led(lexer.MINUS_EQUALS, assignment, parseAssignmentExpression)

	// Logical
	led(lexer.AND, logical, parseBinaryExpression)
	led(lexer.OR, logical, parseBinaryExpression)
	led(lexer.DOT_DOT, logical, parseRangeExpression)

	// Relational
	led(lexer.LESS, relational, parseBinaryExpression)
	led(lexer.LESS_EQUALS, relational, parseBinaryExpression)
	led(lexer.GREATER, relational, parseBinaryExpression)
	led(lexer.GREATER_EQUALS, relational, parseBinaryExpression)
	led(lexer.EQUALS, relational, parseBinaryExpression)
	led(lexer.NOT_EQUALS, relational, parseBinaryExpression)

	// Additive & Multiplicative
	led(lexer.PLUS, additive, parseBinaryExpression)
	led(lexer.DASH, additive, parseBinaryExpression)
	led(lexer.SLASH, multiplicative, parseBinaryExpression)
	led(lexer.STAR, multiplicative, parseBinaryExpression)
	led(lexer.PERCENT, multiplicative, parseBinaryExpression)

	// Literals & Symbols
	nud(lexer.INTEGER, primary, parsePrimaryExpression)
	nud(lexer.FLOAT, primary, parsePrimaryExpression)
	nud(lexer.BOOLEAN, primary, parsePrimaryExpression)
	nud(lexer.STRING, primary, parsePrimaryExpression)
	nud(lexer.IDENTIFIER, primary, parsePrimaryExpression)

	// Unary/Prefix
	nud(lexer.TYPEOF, unary, parsePrefixExpression)
	nud(lexer.DASH, unary, parsePrefixExpression)
	nud(lexer.NOT, unary, parsePrefixExpression)
	nud(lexer.OPEN_BRACKET, primary, parseArrayLiteralExpression)

	// Member / Computed // Call
	led(lexer.DOT, member, parseMemberExpression)
	led(lexer.OPEN_BRACKET, member, parseMemberExpression)
	led(lexer.OPEN_PAREN, call, parseCallExpression)

	// Grouping Expr
	nud(lexer.OPEN_PAREN, defalt_bp, parseGroupingExpression)
	nud(lexer.FUNCTION, defalt_bp, parseFunctionExpression)
	nud(lexer.NEW, defalt_bp, parseNewExpression)

	stmt(lexer.OPEN_CURLY, parseBlockStatement)
	stmt(lexer.LET, parseVariableDeclarationStatement)
	stmt(lexer.VAR, parseVariableDeclarationStatement)
	stmt(lexer.CONST, parseVariableDeclarationStatement)
	stmt(lexer.FUNCTION, parseFunctionDeclaration)
	stmt(lexer.IF, parseIfStatement)
	stmt(lexer.SCOPE, parseScopeStatement)
	stmt(lexer.IMPORT, parseImportStatement)
	stmt(lexer.FOREACH, parseForEachStatement)
	stmt(lexer.CLASS, parseClassDeclarationStatement)
	// TODO: Allow functions to also be declared as public | private | protected | static
	stmt(lexer.PUBLIC, parsePublicScopeDeclarationStatement)
	stmt(lexer.PRIVATE, parsePrivateScopeDeclarationStatement)
	stmt(lexer.PROTECTED, parseProtectedScopeDeclarationStatement)
	stmt(lexer.STATIC, parseStaticDeclarationStatement)
}
