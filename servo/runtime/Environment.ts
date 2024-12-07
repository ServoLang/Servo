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

import { MK_BOOL, MK_NULL, MK_NATIVE_FN, RuntimeValue, MK_NUMBER } from "./Values.ts";

export function createGlobalEnv (): Environment {
  const env = new Environment()
  env.declareVariable("true", MK_BOOL(true), true);
  env.declareVariable("false", MK_BOOL(false), true);
  env.declareVariable("null", MK_NULL(), true);

  // TODO: Do this elsewhere.

  // Define native methods
  env.declareVariable("print", MK_NATIVE_FN((args, scope) => {
    console.log(args[0].value);
    return MK_NULL();
  }), true);

  function timeFunction (args: RuntimeValue[], env: Environment) {
    return MK_NUMBER(Date.now());
  }
  env.declareVariable("time", MK_NATIVE_FN(timeFunction), true);

  return env;
}

export default class Environment {
  private readonly parent?: Environment;
  private variables: Map<string, RuntimeValue>;
  private constants: Set<string>;

  constructor (parentENV?: Environment) {
    this.parent = parentENV;
    this.variables = new Map();
    this.constants = new Set();
  }

  public declareVariable (varname: string, value: RuntimeValue, constant: boolean): RuntimeValue {
    if (this.variables.has(varname)) {
      throw `Cannot declare variable ${varname}. As it already is defined.`;
    }

    this.variables.set(varname, value);
    if (constant) {
      this.constants.add(varname);
    }
    return value;
  }

  public assignVariable (varname: string, value: RuntimeValue): RuntimeValue {
    const env = this.resolve(varname);

    // Cannot assign to constant
    if (env.constants.has(varname)) {
      throw `Cannot reasign to variable ${varname} as it was declared constant.`;
    }

    env.variables.set(varname, value);
    return value;
  }

  public assignAccessor (varname: string, value: RuntimeValue): RuntimeValue {
    const env = this.resolve(varname);
    env.variables.set(varname, value);
    return value;
  }

  public lookupVariable (varname: string): RuntimeValue {
    const env = this.resolve(varname);
    return env.variables.get(varname) as RuntimeValue;
  }

  public resolve (varname: string): Environment {
    if (this.variables.has(varname)) {
      return this;
    }

    if (this.parent == undefined) {
      throw `Cannot resolve '${varname}' as it does not exist.`;
    }

    return this.parent.resolve(varname);
  }
}