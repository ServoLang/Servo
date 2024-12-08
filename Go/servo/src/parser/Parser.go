package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
	"fmt"
	"regexp"
	"strings"
)

type parser struct {
	tokens []lexer.Token
	pos    int
	at     int
	line   int
	file   string
}

func requireFileType(file string) bool {
	if strings.Contains(file, ".servo") || strings.Contains(file, ".svo") {
		return true
	}
	return false
}

// Removes declaration path from file name to only contain the file name
// Irrelevant to reading name, only for error type printing.
func appendFileName(file string) string {
	lastIndex := strings.LastIndex(file, "/")
	if lastIndex != -1 {
		strings.Replace(file, "/"+file[lastIndex+1:], ",", 1)
		file = strings.ReplaceAll(file[lastIndex:], "/", "")
	}
	return file
}

func createParser(tokens []lexer.Token, file string) *parser {
	if !requireFileType(file) {
		panic("file type not supported")
	}

	createTokenLookups()
	createTokenTypeLookups()
	return &parser{tokens: tokens, pos: 0, line: 1, file: appendFileName(file)}
}

func Parse(tokens []lexer.Token, file string) ast.BlockStmt {
	Body := make([]ast.Stmt, 0)
	p := createParser(tokens, file)

	for p.hasTokens() {
		Body = append(Body, parseStatement(p))
	}

	return ast.BlockStmt{Body: Body}
}

// Helper Methods

func (p *parser) getLine() int {
	return p.line
}

func (p *parser) getCharNumber() int {
	return p.at
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) advanceLine() {
	p.line++
	p.at = 1
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	m, _ := regexp.MatchString(`\r?\n`, tk.Value)

	if m {
		p.advanceLine()
	} else {
		p.pos++
		p.at++
	}
	return tk
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf(`Expected %s but got %s instead.\n`, lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind))
		}
		panic(err)
	}

	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}
