import { FunctionDeclaration, Program, VarDeclaration } from "../../frontend/AST.ts";
import Environment from "../Environment.ts";
import { evaluate } from "../Interpreter.ts";
import { FunctionValue, MK_NULL, RuntimeValue } from "../Values.ts";

export function evaluateProgram (program: Program, env: Environment): RuntimeValue {
  let lastEvaluated: RuntimeValue = MK_NULL();
  for (const statement of program.body) {
    lastEvaluated = evaluate(statement, env);
  }
  return lastEvaluated;
}

export function evaluateVariableDeclaration (declaration: VarDeclaration, env: Environment,): RuntimeValue {
  const value = declaration.value
    ? evaluate(declaration.value, env)
    : MK_NULL();

  return env.declareVariable(declaration.identifier, value, declaration.constant);
}

export function evaluateFunctionDeclaration (declaration: FunctionDeclaration, env: Environment,): RuntimeValue {
  const fn = { type: "function", name: declaration.name, parameters: declaration.parameters, declarationEnv: env, body: declaration.body } as FunctionValue;
  return env.declareVariable(declaration.name, fn, true);
}