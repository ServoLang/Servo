// deno-lint-ignore-file
import { FunctionDeclaration, Statement } from "../frontend/AST.ts";
import Environment from "./Environment.ts";

export type ValueType =
    | "null"
    | "number"
    | "boolean"
    | "object"
    | "native"
    | "function"
    | "access";

export interface RuntimeValue {
    type: ValueType;
    value?: any;
}

/**
 * Defines a value of undefined meaning
 */
export interface NullValue extends RuntimeValue {
    type: "null";
    value: null;
}

export function MK_NULL() {
    return { type: "null", value: null } as NullValue;
}

export interface BooleanValue extends RuntimeValue {
    type: "boolean";
    value: boolean | null;
}

export function MK_BOOL(b = true) {
    return { type: "boolean", value: b } as BooleanValue;
}

/**
 * Runtime value that has access to the raw native javascript number.
 */
export interface NumberValue extends RuntimeValue {
    type: "number";
    value: number;
}

export function MK_NUMBER(n = 0) {
    return { type: "number", value: n } as NumberValue;
}

/**
 * Object value
 */
export interface ObjectValue extends RuntimeValue {
    type: "object";
    properties: Map<string, RuntimeValue>;
}

export type FunctionCall = (args: RuntimeValue[], env: Environment) => RuntimeValue;
export interface NativeFnValue extends RuntimeValue {
    type: "native";
    call: FunctionCall;
}

export function MK_NATIVE_FN(call: FunctionCall) {
    return { type: "native", call } as NativeFnValue;
}

export interface FunctionValue extends RuntimeValue {
    type: "function";
    name: string;
    parameters: string[];
    declarationEnv: Environment;
    body: Statement[];
}
