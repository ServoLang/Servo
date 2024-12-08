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

	return parseExpressionStatement(p)
}

func parseExpressionStatement(p *parser) ast.ExpressionStmt {
	expression := parseExpression(p, defalt_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parseBlockStatement(p *parser) ast.Stmt {
	p.expect(lexer.OPEN_CURLY)
	body := []ast.Stmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseStatement(p))
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStmt{
		Body: body,
	}
}

func parseVariableDeclarationStatement(p *parser) ast.Stmt {
	var explicitType ast.Type
	startToken := p.advance().Kind
	isConstant := startToken == lexer.CONST
	symbolName := p.expectError(lexer.IDENTIFIER,
		fmt.Sprintf("Following %s expected variable name however instead recieved %s instead\n",
			lexer.TokenKindString(startToken), lexer.TokenKindString(p.currentTokenKind())))

	if p.currentTokenKind() == lexer.POINTER {
		p.expect(lexer.POINTER)
		explicitType = parseType(p, defalt_bp)
	}

	var assignmentValue ast.Expr
	if p.currentTokenKind() != lexer.POINTER {
		p.expect(lexer.ASSIGNMENT)
		assignmentValue = parseExpression(p, assignment)
	} else if explicitType == nil {
		panic("Missing explicit type for variable declaration.")
	}

	p.expect(lexer.SEMI_COLON)

	if isConstant && assignmentValue == nil {
		panic("Cannot define constant variable without providing default value.")
	}

	return ast.VarDeclarationStmt{
		Constant:      isConstant,
		Identifier:    symbolName.Value,
		AssignedValue: assignmentValue,
		ExplicitType:  explicitType,
	}
}

func parseFunctionParamsAndBody(p *parser) ([]ast.Parameter, ast.Type, []ast.Stmt) {
	functionParams := make([]ast.Parameter, 0)

	p.expect(lexer.OPEN_PAREN)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		// TODO: Flip to state type before identifier
		// Current: function someObject(obj: Object)
		// Needed: function someObject(Object obj)
		paramName := p.expect(lexer.IDENTIFIER).Value
		paramType := parseType(p, defalt_bp)

		functionParams = append(functionParams, ast.Parameter{
			Name: paramName,
			Type: paramType,
		})

		if !p.currentToken().IsOneOfMany(lexer.CLOSE_PAREN, lexer.EOF) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)
	var returnType ast.Type

	if p.currentTokenKind() == lexer.POINTER {
		p.advance()
		returnType = parseType(p, defalt_bp)
	}

	functionBody := ast.ExpectStmt[ast.BlockStmt](parseBlockStatement(p)).Body

	return functionParams, returnType, functionBody
}

func parseFunctionDeclaration(p *parser) ast.Stmt {
	p.advance()
	functionName := p.expect(lexer.IDENTIFIER).Value
	functionParams, returnType, functionBody := parseFunctionParamsAndBody(p)

	return ast.FunctionDeclarationStmt{
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
		Name:       functionName,
	}
}

func parseIfStatement(p *parser) ast.Stmt {
	p.advance()
	condition := parseExpression(p, assignment)
	consequent := parseBlockStatement(p)

	var alternate ast.Stmt
	if p.currentTokenKind() == lexer.ELSE {
		p.advance()

		if p.currentTokenKind() == lexer.IF {
			alternate = parseIfStatement(p)
		} else {
			alternate = parseBlockStatement(p)
		}
	}

	return ast.IfStmt{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parseImportStatement(p *parser) ast.Stmt {
	p.advance()
	var importFrom string
	importName := p.expect(lexer.IDENTIFIER).Value

	if p.currentTokenKind() == lexer.FROM {
		p.advance()
		importFrom = p.expect(lexer.STRING).Value
	} else {
		importFrom = importName
	}

	p.expect(lexer.SEMI_COLON)
	return ast.ImportStmt{
		Name: importName,
		From: importFrom,
	}
}

func parseForEachStatement(p *parser) ast.Stmt {
	p.advance()
	p.expect(lexer.OPEN_PAREN)
	valueName := p.expect(lexer.IDENTIFIER).Value

	var index bool
	if p.currentTokenKind() == lexer.COMMA {
		p.expect(lexer.COMMA)
		p.expect(lexer.IDENTIFIER)
		index = true
	}

	p.expect(lexer.IN)
	iterable := parseExpression(p, defalt_bp)
	p.expect(lexer.CLOSE_PAREN)
	body := ast.ExpectStmt[ast.BlockStmt](parseBlockStatement(p)).Body

	return ast.ForeachStmt{
		Value:    valueName,
		Index:    index,
		Iterable: iterable,
		Body:     body,
	}
}

func parseClassDeclarationStatement(p *parser) ast.Stmt {
	p.advance()
	className := p.expect(lexer.IDENTIFIER).Value
	classBody := parseBlockStatement(p)

	return ast.ClassDeclarationStmt{
		Name: className,
		Body: ast.ExpectStmt[ast.BlockStmt](classBody).Body,
	}
}
