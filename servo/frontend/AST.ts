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

// -----------------------------------------------------------
// --------------          AST TYPES        ------------------
// ---     Defines the structure of our languages AST      ---
// -----------------------------------------------------------

export type NodeType =
  // STATEMENTS
  | "Program"
  | "VarDeclaration"
  | "FunctionDeclaration"
  | "AccessDeclaration"
  // EXPRESSIONS
  | "AssignmentExpr"
  | "MemberExpr"
  | "CallExpr"

  // LITERALS
  | "Property"
  | "ObjectLiteral"
  | "NumericLiteral"
  | "Identifier"
  | "BinaryExpr";

/**
 * Statements do not result in a value at runtime.
 They contain one or more expressions internally */
export interface Statement {
  kind: NodeType;
}

/**
 * Defines a block which contains many statements.
 * -  Only one program will be contained in a file.
 */
export interface Program extends Statement {
  kind: "Program";
  body: Statement[];
}

export interface VarDeclaration extends Statement {
  kind: "VarDeclaration";
  constant: boolean;
  identifier: string;
  value?: Expression;
}

export interface FunctionDeclaration extends Statement {
  kind: "FunctionDeclaration";
  parameters: string[]; // THIS PREVENTS NEEDING TO GIVE PARAMETERS A VALUE IN A METHOD
  name: string;
  body: Statement[];
  protected?: boolean;
  private?: boolean;
  public?: boolean;
  async?: boolean;
}

/**
 * Access Declaration public | private | protected.
 */
export interface AccessDeclaration extends Statement {
    kind: "AccessDeclaration";
    name: string;
    function: FunctionDeclaration;
    protected?: boolean;
    private?: boolean;
    public?: boolean;
}

/**  Expressions will result in a value at runtime unlike Statements */
export interface Expression extends Statement {}

export interface AssignmentExpr extends Expression {
  kind: "AssignmentExpr";
  assigne: Expression;
  value: Expression;
}

/**
 * An operation with two sides seperated by an operator.
 * Both sides can be ANY Complex Expression.
 * - Supported Operators -> + | - | / | * | % | ^
 */
export interface BinaryExpr extends Expression {
  kind: "BinaryExpr";
  left: Expression;
  right: Expression;
  operator: string; // needs to be of type BinaryOperator
}

export interface CallExpr extends Expression {
  kind: "CallExpr";
  args: Expression[];
  caller: Expression;
}

export interface MemberExpr extends Expression {
  kind: "MemberExpr";
  object: Expression;
  property: Expression;
  computed: boolean;
}

// LITERAL / PRIMARY EXPRESSION TYPES
/**
 * Represents a user-defined variable or symbol in source.
 */
export interface Identifier extends Expression {
  kind: "Identifier";
  symbol: string;
}

/**
 * Represents a numeric constant inside the source code.
 */
export interface NumericLiteral extends Expression {
  kind: "NumericLiteral";
  value: number;
}

export interface Property extends Expression {
  kind: "Property";
  key: string,
  value?: Expression,
}

export interface ObjectLiteral extends Expression {
  kind: "ObjectLiteral";
  properties: Property[];
}