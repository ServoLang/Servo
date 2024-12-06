import { AssignmentExpr, BinaryExpr, CallExpr, FunctionDeclaration, Identifier, NumericLiteral, ObjectLiteral, Program, Statement, VarDeclaration } from "../frontend/AST.ts";
import Environment from "./Environment.ts";
import { evaluateFunctionDeclaration, evaluateProgram, evaluateVariableDeclaration } from "./eval/Statements.ts";
import { evaluateAssignment, evaluateBinaryExpression, evaluateCallExpression, evaluateIdentifier, evaluateObjectExpression } from "./eval/Expressions.ts";
import { NumberValue, RuntimeValue } from "./Values.ts";

export function evaluate (astNode: Statement, env: Environment): RuntimeValue {
    switch (astNode.kind) {
        case "NumericLiteral":
            return {
                value: ((astNode as NumericLiteral).value),
                type: "number",
            } as NumberValue;
        case "Identifier":
            return evaluateIdentifier(astNode as Identifier, env);
        case "ObjectLiteral":
            return evaluateObjectExpression(astNode as ObjectLiteral, env);
        case "CallExpr":
            return evaluateCallExpression(astNode as CallExpr, env);
        case "AssignmentExpr":
            return evaluateAssignment(astNode as AssignmentExpr, env);
        case "BinaryExpr":
            return evaluateBinaryExpression(astNode as BinaryExpr, env);
        case "Program":
            return evaluateProgram(astNode as Program, env);
        // Handle statements
        case "VarDeclaration":
            return evaluateVariableDeclaration(astNode as VarDeclaration, env);
        case "FunctionDeclaration":
            return evaluateFunctionDeclaration(astNode as FunctionDeclaration, env);
            // Handle unimplemented ast types as error.
        default:
            console.error("This AST Node has not yet been setup for interpretation.", astNode);
            Deno.exit(0);
    }
}
