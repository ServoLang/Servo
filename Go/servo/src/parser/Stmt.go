package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
)

func parse_stmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)
	return ast.ExpressionStmt{Expression: expression}
}

func parse_var_decl_stmt(p *parser) ast.Stmt {
	var explicitType ast.Type
	var assignedValue ast.Expr
	isConstant := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name.").Value

	// Explicit type could be present
	if p.currentTokenKind() == lexer.POINTER {
		p.advance() // eat
		explicitType = parse_type(p, default_bp)
	}

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignedValue = parse_expr(p, assignment)
	} else if explicitType == nil {
		panic(`Missing either type declaration or value. Expected to find a variable declaration.`)
	}

	p.expect(lexer.SEMI_COLON)

	if isConstant && assignedValue == nil {
		panic(`Cannot define constant without providing a value.`)
	}

	return ast.VarDeclStmt{ExplicitType: explicitType, IsConstant: isConstant, VariableName: varName, AssignedValue: assignedValue}
}
