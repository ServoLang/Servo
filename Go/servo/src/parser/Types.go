package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
	"fmt"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp BindingPower) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]BindingPower

// Lookup tables
var type_bp_lu = type_bp_lookup{}
var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}

func typeLED(kind lexer.TokenKind, bp BindingPower, led_fn type_led_handler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_fn
}

func typeNUD(kind lexer.TokenKind, nud_fn type_nud_handler) {
	type_bp_lu[kind] = primary
	type_nud_lu[kind] = nud_fn
}

func createTokenTypeLookups() {
	typeNUD(lexer.IDENTIFIER, parseSymbolType)
	typeLED(lexer.OPEN_BRACKET, call, parseArrayType)
}

func parseSymbolType(p *parser) ast.Type {
	return ast.SymbolType{Name: p.expect(lexer.IDENTIFIER).Value}
}

// Syntax for let obj -> Object[][] = SomeObject;
func parseArrayType(p *parser, left ast.Type, bp BindingPower) ast.Type {
	p.expect(lexer.OPEN_BRACKET)
	p.expect(lexer.CLOSE_BRACKET)
	return ast.ArrayType{Underlying: left}
}

func parseType(p *parser, bp BindingPower) ast.Type {
	// Parse NUD
	tokenKind := p.currentTokenKind()
	nudFunction, exists := type_nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token '%s' | line:%d:%d ('%s')\n", lexer.TokenKindString(tokenKind), p.getLine(), p.getCharNumber(), p.file))
	}

	left := nudFunction(p)
	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFunction, exists := type_led_lu[tokenKind]
		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token '%s' | line:%d:%d ('%s')\n", lexer.TokenKindString(tokenKind), p.getLine(), p.getCharNumber(), p.file))
		}

		left = ledFunction(p, left, bp)
	}
	return left
}
