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
// ---------------          LEXER          -------------------
// ---  Responsible for producing tokens from the source   ---
// -----------------------------------------------------------

// Represents tokens that our language understands in parsing.
import { ErrorType, Exception } from "../runtime/error/Exception.ts";

export enum TokenType {
  // Literal Types
  Number,
  Identifier,
  // Keywords
  Var,
  Let,
  Const,
  Num,
  String,
  New,
  Throw,
  Optional,

  Scope,
  Class,
  With,
  Super,

  Public,
  Private,
  Protected,
  Function,
  Void, // empty return type function
  Return,
  Async,
  Static,
  Arrow,

  // Turing Operators
  If,
  Then,
  Else,
  For,
  While,
  Do,
  Break,

  // Grouping * Operators
  BinaryOperator,
  Equals,
  Comma,
  Dot,
  Colon,
  Semicolon,
  OpenParen,
  CloseParen,
  OpenBrace,
  CloseBrace,
  OpenBracket,
  CloseBracket,
  LessThan,
  GreaterThan,
  DoubleQuote,
  Quote,
  Grave,
  Tilde,
  Exclamation,
  EOF, // Signified the end of file
  Unknown,
}

/**
 * Constant lookup for keywords and known identifiers + symbols.
 */
const KEYWORDS: Record<string, TokenType> = {
  let: TokenType.Let,
  const: TokenType.Const,
  var: TokenType.Var,
  num: TokenType.Num, // TODO: not fully implemented
  string: TokenType.String,  // unimplemented
  new: TokenType.New, // unimplemented
  throw: TokenType.Throw, // unimplemented
  optional: TokenType.Optional,  // unimplemented
  arrow: TokenType.Arrow,

  scope: TokenType.Scope, // unimplemented
  class: TokenType.Class, // unimplemented
  with: TokenType.With, // unimplemented
  super: TokenType.Super, // unimplemented

  public: TokenType.Public, // unimplemented
  private: TokenType.Private, // unimplemented
  protected: TokenType.Protected, // unimplemented
  function: TokenType.Function,
  Void: TokenType.Void,
  return: TokenType.Return, // unimplemented
  async: TokenType.Async,  // unimplemented
  static: TokenType.Static, // unimplemented

  if: TokenType.If, // unimplemented
  then: TokenType.Then, // unimplemented
  else: TokenType.Else, // unimplemented
  for: TokenType.For, // unimplemented
  while: TokenType.While, // unimplemented
  do: TokenType.Do, // unimplemented
  break: TokenType.Break, // unimplemented
};

// Represents a single token from the source-code.
export interface Token {
  value: string; // contains the raw value as seen inside the source code.
  type: TokenType; // tagged structure.
}

// Returns a token of a given type and value
function token (value = "", type: TokenType): Token {
  return { value, type };
}

/**
 * Returns whether the character passed in alphabetic -> [a-zA-Z]
 */
function isAlpha (src: string) {
  return src.toUpperCase() != src.toLowerCase();
}

/**
 * Returns true if the character is whitespace like -> [\s, \t, \n]
 */
function isSkippable (str: string) {
  return str == " " || str == "\n" || str == "\t" || str == "\r";
}

/**
 * Returns true if either of the characters is an arrow type.
 */
function isArrow (str: string) {
  return str == "-" || str == ">";
}

/**
 Return whether the character is a valid integer -> [0-9]
 */
function isInt (str: string) {
  const c = str.charCodeAt(0);
  const bounds = ["0".charCodeAt(0), "9".charCodeAt(0)];
  return c >= bounds[0] && c <= bounds[1];
}

function tokenFor (str: string): TokenType {
  let c: TokenType = TokenType.Unknown;
  switch (str) {
    case '"': {
      c = TokenType.DoubleQuote;
      break;
    }
    case "`": {
      c = TokenType.Grave;
      break;
    }
    case "'": {
      c = TokenType.Quote;
      break;
    }
  }
    return c;
}

function pushString (quote: string, src: string[], tokens: Token[]): Token {
  console.log(tokens);
  let str = "";
  src.shift(); // Skip the opening quote.
  while (src.length > 0 && src[0] != quote) {
    str += src.shift();
  }

  console.log(`String: ${str}`);
  let temp = "";
  temp += src.shift();// Skip the closing quote.
  tokens.push(token(temp, tokenFor(temp)));
  return token(str, TokenType.String);
}

/**
 * Given a string representing source code: Produce tokens and handles
 * possible unidentified characters.
 *
 * - Returns a array of tokens.
 * - Does not modify the incoming string.
 */
export function tokenize (sourceCode: string): Token[] {
  const tokens = new Array<Token>();
  const src = sourceCode.split("");

  // produce tokens until the EOF is reached.
  while (src.length > 0) {
    // BEGIN PARSING ONE CHARACTER TOKENS

    // Handle Arrow Tokens
    if (src[0] == "-" && src[1] == ">") {
      let arrow = "";
      while (src.length > 0 && isArrow(src[0])) {
        arrow += src.shift();
      }
      tokens.push(token(arrow, TokenType.Arrow));
    } else if (src[0] == "(") {
      tokens.push(token(src.shift(), TokenType.OpenParen));
    } else if (src[0] == ")") {
      tokens.push(token(src.shift(), TokenType.CloseParen));
    } else if (src[0] == "{") {
      tokens.push(token(src.shift(), TokenType.OpenBrace));
    } else if (src[0] == "}") {
      tokens.push(token(src.shift(), TokenType.CloseBrace));
    } else if (src[0] == "[") {
      tokens.push(token(src.shift(), TokenType.OpenBracket));
    } else if (src[0] == "]") {
      tokens.push(token(src.shift(), TokenType.CloseBracket));
    } // HANDLE BINARY OPERATORS
    else if (src[0] == "+" || src[0] == "-" || src[0] == "*" || src[0] == "/" || src[0] == "%" || src[0] == "^") {
      tokens.push(token(src.shift(), TokenType.BinaryOperator));
    } // Handle Conditional & Assignment Tokens
    else if (src[0] == "=") {
      tokens.push(token(src.shift(), TokenType.Equals));
    } else if (src[0] == ";") {
      tokens.push(token(src.shift(), TokenType.Semicolon));
    } else if (src[0] == ":") {
      tokens.push(token(src.shift(), TokenType.Colon));
    } else if (src[0] == ",") {
      tokens.push(token(src.shift(), TokenType.Comma));
    } else if (src[0] == ".") {
      tokens.push(token(src.shift(), TokenType.Dot));
    } else if (src[0] == "<") {
      tokens.push(token(src.shift(), TokenType.LessThan));
    } else if (src[0] == ">") {
      tokens.push(token(src.shift(), TokenType.GreaterThan));
    } else if (src[0] == "\"" || src[0] == "'" || src[0] == "`") {
      tokens.push(token(src[0], tokenFor(src[0])));
      if (src[1] != ";") tokens.push(pushString(src[0], src, tokens));
    } else if (src[0] == "~") {
      tokens.push(token(src.shift(), TokenType.Tilde));
    } else if (src[0] == "!") {
      tokens.push(token(src.shift(), TokenType.Exclamation));
    } // HANDLE MULTICHARACTER KEYWORDS, TOKENS, IDENTIFIERS ETC...
    else {
      // Handle numeric literals -> Integers
      if (isInt(src[0])) {
        let num = "";
        while (src.length > 0 && isInt(src[0])) {
          num += src.shift();
        }

        // append new numeric token.
        tokens.push(token(num, TokenType.Number));
      }
      // Handle Identifier & Keyword Tokens.
      else if (isAlpha(src[0])) {
        let ident = "";
        while (src.length > 0 && isAlpha(src[0])) {
          ident += src.shift();
        }

        // CHECK FOR RESERVED KEYWORDS
        const reserved = KEYWORDS[ident];
        // If value is not undefined then the identifier is
        // recognized keyword
        if (typeof reserved == "number") {
          tokens.push(token(ident, reserved));
        } else {
          // Unrecognized name must mean user defined symbol.
          tokens.push(token(ident, TokenType.Identifier));
        }
      } else if (isSkippable(src[0])) {
        // Skip unneeded chars.
        src.shift();
      } // Handle unrecognized characters.
      // TODO: Implement better errors and error recovery.
      else {
        throw new Exception(ErrorType.InternalError, `Unrecognized character found in source: char#:${src[0].charCodeAt(0)} : '${src[0]}'`).exit();
      }
    }
  }

  tokens.push({ type: TokenType.EOF, value: "EndOfFile" });
  return tokens;
}