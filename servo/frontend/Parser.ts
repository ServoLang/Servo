/*
 * Copyright (c) 2024. Servo Contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

// deno-lint-ignore-file no-explicit-any
import { AccessDeclaration, AssignmentExpr, BinaryExpr, CallExpr, Expression, FunctionDeclaration, Identifier, MemberExpr, NumericLiteral, ObjectLiteral, Program, Property, Statement, StringLiteral, VarDeclaration } from "./AST.ts";

import { Token, tokenize, TokenType } from "./Lexer.ts";

/**
 * Frontend for producing a valid AST (Abstract Syntax Tree) from sourcecode
 */
export default class Parser {
	private tokens: Token[] = [];

	/**
	 * Determines if the parsing is complete and the END OF FILE Is reached.
	 */
	private notEOF (): boolean {
		return this.tokens[0].type != TokenType.EOF;
	}

	/**
	 * Returns the currently available token
	 */
	private at () {
		return this.tokens[0] as Token;
	}

	/**
	 * Returns the previous token and then advances the tokens array to the next value.
	 */
	private eat () {
		return this.tokens.shift() as Token;
	}

	/**
	 * Returns the previous token and then advances the tokens array to the next value.
	 *  Also checks the type of expected token and throws if the values do not match.
	 */
	private expect (type: TokenType, err: any) {
		const prev = this.tokens.shift() as Token;
		if (!prev || prev.type != type) {
			console.error("Parser Error:\n", err, prev, " - Expecting: ", type);
			Deno.exit(1);
		}

		return prev;
	}

	public produceAST (sourceCode: string): Program {
		this.tokens = tokenize(sourceCode);
		const program: Program = { kind: "Program", body: [], };

		// Parse until end of file
		while (this.notEOF()) {
			program.body.push(this.parseStatement());
		}

		return program;
	}

	// Handle complex statement types
	private parseStatement (): Statement {
		// skip to parse_expr
		switch (this.at().type) {
			case TokenType.Var:
			case TokenType.Let:
			case TokenType.Const:
				return this.parseVariableDeclaration();
			case TokenType.Num:
				return this.parseNumberDeclaration();
			//case TokenType.String:
				//return this.parseStringDeclaration();
			case TokenType.Public:
			case TokenType.Private:
			case TokenType.Protected:
				return this.parseAccessDeclaration();
			case TokenType.Function:
				return this.parseFunctionDeclaration();
			default:
				return this.parseExpression();
		}
	}

	private parseAccessDeclaration (): Statement {
		this.eat(); // eat access keyword
		const func = this.parseFunctionDeclaration();
		return { kind: "AccessDeclaration", function: func } as AccessDeclaration;
	}

	private parseFunctionDeclaration (): Statement {
		this.eat(); // eat declaration keyword
		const name = this.expect(TokenType.Identifier,"Expected function name following declaration type keyword.").value;
		const args = this.parseArguments();
		const params: string[] = [];

		for (const arg of args) {
			if (arg.kind !== "Identifier") {
				console.log(arg);
				throw `Inside function declaration expected parameters to be of type string.`;
			}

			params.push((arg as Identifier).symbol);
		}

		this.expect(TokenType.Arrow, "Expected arrow token following function parameters.");
		// TODO: Change to multiple return types. Currently only supports voids.
		this.expect(TokenType.Void, "Expected return type for function declaration.");
		this.expect(TokenType.OpenBrace,"Expected function body following function declaration.");

		const body: Statement[] = [];

		while (this.at().type !== TokenType.EOF && this.at().type !== TokenType.CloseBrace) {
			body.push(this.parseStatement());
		}

		this.expect(TokenType.CloseBrace,"Expected closure of function declaration.");
		return { body, name, parameters: params, kind: "FunctionDeclaration" } as FunctionDeclaration;
	}

	// LET IDENT;
	// ( LET | CONST | VAR) IDENT = EXPR;
	private parseVariableDeclaration(): Statement {
		const isConstant = this.eat().type == TokenType.Const;
		const identifier = this.expect(TokenType.Identifier,"Expected identifier name following let | const | var keywords.").value;

		if (this.at().type == TokenType.Semicolon) {
			this.eat(); // expect semicolon
			if (isConstant) {
				throw "Must assign value to constant expression. No value provided.";
			}

			this.expect(TokenType.Semicolon, "Variable declaration statement must end with semicolon.");
			return { kind: "VarDeclaration", identifier, constant: false } as VarDeclaration;
		}

		this.expect(TokenType.Equals,"Expected equals token following identifier in var declaration.",);
		const declaration = { kind: "VarDeclaration", value: this.parseExpression(), identifier, constant: isConstant } as VarDeclaration;
		this.expect(TokenType.Semicolon, "Variable declaration statement must end with semicolon.");

		return declaration;
	}

	// TODO: PROBABLY DOESN'T WORK AS INTENDED
	// SET TO ONLY ALLOW NUMERIC LITERALS AFTER EQUALS
	private parseNumberDeclaration(): Statement {
		const isNum = this.eat().type == TokenType.Num;
		const identifier = this.expect(TokenType.Identifier, "Expected identifier name following num keyword.",).value;

		if (this.at().type == TokenType.Semicolon) {
			this.eat();
			return { kind: "VarDeclaration", identifier, constant: false, } as VarDeclaration;
		}

		this.expect(TokenType.Equals,"Expected equals token following identifier in num declaration.");
		const declaration = { kind: "VarDeclaration", value: this.parseExpression(), identifier, constant: false } as VarDeclaration;
		this.expect(TokenType.Semicolon, "Variable declaration statement must end with semicolon.");

		return declaration;
	}

	// TODO: Properly implement string declaration. Needs to be able to handle string literals.
	private parseStringDeclaration(): Statement {
		const isString = this.eat().type == TokenType.String;
		const identifier = this.expect(TokenType.Identifier,"Expected identifier name following num keyword.").value;

		if (this.at().type == TokenType.Semicolon) {
			this.eat();
			return { kind: "VarDeclaration", identifier, constant: false, } as VarDeclaration;
		}

		this.expect(TokenType.Equals, "Expected equals token following identifier in num declaration.",);
		const beginningToken = this.expect(TokenType.Tilde | TokenType.Quote | TokenType.DoubleQuote, "Expected a string value following equals token.").value;
		const declaration = { kind: "VarDeclaration", value: this.parseExpression(), identifier, constant: false } as VarDeclaration;
		const endingToken = this.expect(TokenType.Tilde | TokenType.Quote | TokenType.DoubleQuote, "Expected a string value following equals token.").value;
		this.expect(TokenType.Semicolon,"Variable declaration statement must end with semicolon.");

		if (beginningToken != endingToken) {
			throw `String declaration must be closed with the same token as it was opened with. Expected ${beginningToken} but found ${endingToken}.`;
		}

		return declaration;
	}

	// Handle expressions
	private parseExpression(): Expression {
		return this.parseAssignmentExpression();
	}

	private parseAssignmentExpression(): Expression {
		const left = this.parseObjectExpression();

		if (this.at().type == TokenType.Equals) {
			this.eat(); // advance past equals
			const value = this.parseAssignmentExpression();
			this.expect(TokenType.Semicolon, "Variable declaration statement must end with semicolon.");
			return { value, assigne: left, kind: "AssignmentExpr" } as AssignmentExpr;
		}

		return left;
	}

	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// TODO: Will need to reformat to make into using parenthesis for objects rather than open braces
	// Open braces for objects is ugly af
	private parseObjectExpression(): Expression {
		if (this.at().type !== TokenType.OpenBrace) {
			return this.parseAdditiveExpression();
		}

		this.eat(); // advance
		const properties = new Array<Property>();

		while (this.notEOF() && this.at().type != TokenType.CloseBrace) {
			// {key: val, key2: val}
			const key = this.expect(TokenType.Identifier,"Object literal key expected.").value;

			// allows shorthand key: pair -> { key, }
			if (this.at().type == TokenType.Comma) {
				this.eat(); // advance past a comma
				properties.push({ key, kind: "Property", value: undefined } as Property);
				continue;
			} // allows shorthand key: pair -> { key }
			else if (this.at().type == TokenType.CloseBrace) {
				properties.push({ key, kind: "Property", value: undefined });
				continue;
			}
			// { key: val }
			this.expect(TokenType.Colon, "Missing colon following identifier in ObjectExpr.");
			const value = this.parseExpression();

			properties.push({ kind: "Property", value, key });
			if (this.at().type != TokenType.CloseBrace) {
				this.expect(TokenType.Comma, "Expected comma or closing bracket following property.");
			}
		}

		this.expect(TokenType.CloseBrace, "Object literal missing closure.");
		this.expect(TokenType.Semicolon, "Variable declaration statement must end with semicolon.");
		return { kind: "ObjectLiteral", properties } as ObjectLiteral;
	}

	// Handle Addition & Subtraction Operations
	private parseAdditiveExpression(): Expression {
		let left = this.parseMultiplicativeExpression();

		while (this.at().value == "+" || this.at().value == "-") {
			const operator = this.eat().value;
			const right = this.parseMultiplicativeExpression();
			left = { kind: "BinaryExpr", left, right, operator } as BinaryExpr;
		}

		return left;
	}

	// Handle Multiplication, Division, Power of & Modulo Operations
	private parseMultiplicativeExpression(): Expression {
		let left = this.parseCallMemberExpression();

		while (this.at().value == "/" || this.at().value == "*" || this.at().value == "%" || this.at().value == "^") {
			const operator = this.eat().value;
			const right = this.parseCallMemberExpression();
			left = { kind: "BinaryExpr", left, right, operator } as BinaryExpr;
		}

		return left;
	}

	private parseCallMemberExpression(): Expression {
		const member = this.parseMemberExpression();

		if (this.at().type == TokenType.OpenParen) {
			return this.parseCallExpression(member);
		}

		return member;
	}

	private parseCallExpression(caller: Expression): Expression {
		let call_expr: Expression = { kind: "CallExpr", caller, args: this.parseArguments() } as CallExpr;

		if (this.at().type == TokenType.OpenParen) {
			call_expr = this.parseCallExpression(call_expr);
		}

		return call_expr;
	}

	private parseArguments(): Expression[] {
		this.expect(TokenType.OpenParen, "Expected open parenthesis."); // Error should never appear.
		const args = this.at().type == TokenType.CloseParen ? [] : this.parseArgumentsList();

		this.expect(TokenType.CloseParen,"Missing closing parenthesis inside arguments list.");
		return args;
	}

	// TODO: Handle types of arguments instead of undefined list.
	private parseArgumentsList(): Expression[] {
		const args = [this.parseAssignmentExpression()];

		while (this.at().type == TokenType.Comma && this.eat()) {
			args.push(this.parseAssignmentExpression());
		}

		return args;
	}

	// foo.x()()
	private parseMemberExpression(): Expression {
		let object = this.parsePrimaryExpression();

		while (this.at().type == TokenType.Dot || this.at().type == TokenType.OpenBracket) {
			const operator = this.eat();
			let property: Expression;
			let computed: boolean;

			// non computed values aka dot.expr
			if (operator.type == TokenType.Dot) {
				computed = false;
				// get identifier
				property = this.parsePrimaryExpression();

				if (property.kind != "Identifier") {
					throw `Cannot use dot operator without right hand side being an identifier.`;
				}
			} else { // allows obj[computedValue]
				computed = true;
				property = this.parseExpression();
				this.expect(TokenType.CloseBracket,"Missing closing bracket in computed value.",);
			}

			object = { kind: "MemberExpr", object, property, computed } as MemberExpr;
		}

		return object;
	}

	// Orders Of Presidency
	// Assignment
	// Object
	// Additive Expression
	// Multiplicative Expression
	// Call
	// Member
	// Primary Expression

	// Parse Literal Values & Grouping Expressions
	private parsePrimaryExpression(): Expression {
		const tk = this.at().type;

		// Determine which token we are currently at and return literal value
		switch (tk) {
			// User defined values.
			case TokenType.Identifier: {
				return { kind: "Identifier", symbol: this.eat().value } as Identifier;
			}

			// Constants and Numeric Constants
			case TokenType.Number: {
				return { kind: "NumericLiteral", value: parseFloat(this.eat().value) } as NumericLiteral;
			}

			case TokenType.String: {
				return { kind: "StringLiteral", value: this.eat().value } as StringLiteral;
			}

			// Grouping Expressions
			case TokenType.OpenParen: {
				this.eat(); // eat the opening paren
				const value = this.parseExpression();
				this.expect(TokenType.CloseParen,"Unexpected token found inside parenthesised expression. Expected closing parenthesis."); // closing paren
				return value;
			}

			// Unidentified Tokens and Invalid Code Reached
			default:
				console.error("Unexpected token found during parsing!", this.at());
				Deno.exit(1);
		}
	}
}
