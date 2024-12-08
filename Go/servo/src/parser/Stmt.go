package parser

import (
	"Go/servo/src/ast"
	"Go/servo/src/lexer"
	"fmt"
)

func parseStatement(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parseExpression(p, default_bp)
	p.expect(lexer.SEMI_COLON)
	return ast.ExpressionStmt{Expression: expression}
}

func parseVariableDeclarationStatement(p *parser) ast.Stmt {
	var explicitType ast.Type
	isConstant := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name.").Value

	// Explicit type could be present
	if p.currentTokenKind() == lexer.POINTER {
		p.advance() // eat
		explicitType = parseType(p, default_bp)
	}

	var assignedValue ast.Expression
	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignedValue = parseExpression(p, assignment)
	} else if explicitType == nil {
		err := fmt.Sprintf(`Missing either type declaration or value. Expected to find a variable declaration.`)
		panic(err)
	}

	p.expect(lexer.SEMI_COLON)

	if isConstant && assignedValue == nil {
		panic(`Cannot define constant without providing a value.`)
	}

	return ast.VarDeclStmt{ExplicitType: explicitType, IsConstant: isConstant, VariableName: varName, AssignedValue: assignedValue}
}

// FOR CREATING STRUCTS
func parseStructDeclarationStatement(p *parser) ast.Stmt {
	p.expect(lexer.STRUCT) // advance past struct keyword
	var properties = map[string]ast.StructProperty{}
	var methods = map[string]ast.StructMethod{}
	var structName = p.expect(lexer.IDENTIFIER).Value // Struct name

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var isStatic bool
		var propertyName string

		// TODO: Add Public, Private, Protected here

		if p.currentTokenKind() == lexer.STATIC {
			isStatic = true
			p.expect(lexer.STATIC)
		}

		// Property
		if p.currentTokenKind() == lexer.IDENTIFIER {
			propertyName = p.expect(lexer.IDENTIFIER).Value
			p.expectError(lexer.POINTER, "Expected to find a pointer following a property declaration.")
			structType := parseType(p, default_bp)
			p.expect(lexer.SEMI_COLON)

			_, exists := properties[propertyName]
			if exists {
				panic(fmt.Sprintf("Property %s already defined.", propertyName))
			}

			properties[propertyName] = ast.StructProperty{IsStatic: isStatic, Type: structType}
			continue
		}

		// TODO: Allow handling methods
		panic("Cannot currently handle methods inside struct declaration.")

	}

	p.expect(lexer.CLOSE_CURLY)
	// PUT SEMI COLON HERE TO REQUIRE SEMI COLON AFTER CLOSING CURLY
	return ast.StructDeclStmt{Properties: properties, Methods: methods, StructName: structName}
}
