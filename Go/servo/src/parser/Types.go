package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
	"fmt"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp binding_power) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]binding_power

var type_bp_lu = type_bp_lookup{}
var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}

func typeLed(kind lexer.TokenKind, bp binding_power, led_fn type_led_handler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_fn
}

func typeNud(kind lexer.TokenKind, bp binding_power, nud_fn type_nud_handler) {
	type_bp_lu[kind] = primary
	type_nud_lu[kind] = nud_fn
}

func createTypeTokenLookups() {

	typeNud(lexer.IDENTIFIER, primary, func(p *parser) ast.Type {
		return ast.SymbolType{
			Value: p.advance().Value,
		}
	})

	// []number
	typeNud(lexer.OPEN_BRACKET, member, func(p *parser) ast.Type {
		p.advance()
		p.expect(lexer.CLOSE_BRACKET)
		insideType := parseType(p, defalt_bp)

		return ast.ListType{
			Underlying: insideType,
		}
	})
}

func parseType(p *parser, bp binding_power) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}
