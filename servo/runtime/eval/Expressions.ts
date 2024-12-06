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

import { AssignmentExpr, BinaryExpr, CallExpr, Identifier, ObjectLiteral } from "../../frontend/AST.ts";
import Environment from "../Environment.ts";
import { evaluate } from "../Interpreter.ts";
import { FunctionValue, MK_NULL, NativeFnValue, NumberValue, ObjectValue, RuntimeValue } from "../Values.ts";
import { ErrorType, Exception } from "../error/Exception.ts";

// DO BINARY MATHEMATICS OPERATIONS
function evaluateNumericBinaryExpression (lhs: NumberValue, rhs: NumberValue, operator: string): NumberValue {
  let result: number;
  if (operator == "+") {
    result = lhs.value + rhs.value;
  } else if (operator == "-") {
    result = lhs.value - rhs.value;
  } else if (operator == "*") {
    result = lhs.value * rhs.value;
  } else if (operator == "/") {
    if (lhs.value == 0 || rhs.value == 0) throw new Exception(ErrorType.DivideByZero, "Division by zero");
    result = lhs.value / rhs.value;
  } else if (operator == "^") {
    result = Math.pow(lhs.value, rhs.value);
  } else {
    result = lhs.value % rhs.value;
  }

  return { value: result, type: "number" };
}

/**
 * Evaluates expressions following the binary operation type.
 */
export function evaluateBinaryExpression (binop: BinaryExpr, env: Environment,): RuntimeValue {
  const lhs = evaluate(binop.left, env);
  const rhs = evaluate(binop.right, env);

  // Only currently support numeric operations
  if (lhs.type == "number" && rhs.type == "number") {
    return evaluateNumericBinaryExpression(lhs as NumberValue, rhs as NumberValue, binop.operator,);
  }

  // One or both are NULL
  return MK_NULL();
}

export function evaluateIdentifier (ident: Identifier, env: Environment,): RuntimeValue {
  const val = env.lookupVariable(ident.symbol);
  return val;
}

export function evaluateAssignment (node: AssignmentExpr, env: Environment,): RuntimeValue {
  if (node.assigne.kind !== "Identifier") {
    throw new Exception(ErrorType.InvalidAssignment,`Invalid LHS inside assignment expr ${JSON.stringify(node.assigne)}`).exit();
  }

  const varname = (node.assigne as Identifier).symbol;
  return env.assignVariable(varname, evaluate(node.value, env));
}

export function evaluateObjectExpression (obj: ObjectLiteral, env: Environment): RuntimeValue {
  const object = { type: "object", properties: new Map() } as ObjectValue;

  for (const { key, value } of obj.properties) {
    // handles valid key: pair
    const runtimeVal = (value == undefined) ? env.lookupVariable(key) : evaluate(value, env);
    object.properties.set(key, runtimeVal);
  }

  return object;
}

export function evaluateCallExpression (expr: CallExpr, env: Environment): RuntimeValue {
  const args = expr.args.map((arg) => evaluate(arg, env));
  const fn = evaluate(expr.caller, env);

  if (fn.type == "native") {
    const result = (fn as NativeFnValue).call(args, env);
    return result;
  } 
  
  if (fn.type == "function") {
    const func = fn as FunctionValue;
    const scope = new Environment(func.declarationEnv);

    // create variables for the parameters list
    for (let i = 0; i < func.parameters.length; i++) {
      // TODO: check the bounds here.
      // verify arity of function.
      const varname = func.parameters[i];
      scope.declareVariable(varname, args[i], false); // false = prevent assignment into a constructor
    }

    let result: RuntimeValue = MK_NULL();

    // evaluate the function body line by line
    for (const stmt of func.body) {
      result = evaluate(stmt, scope);
    }

    return result;
  }

  throw `Cannot call value that is not a function: ${JSON.stringify(fn)}`;

}
