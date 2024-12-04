import { MK_BOOL, MK_NULL, RuntimeVal } from "./values.ts";

export default class Environment {
    private parent?: Environment;
    private variables: Map<string, RuntimeVal>;
    private constants: Set<string>;

    constructor (parentENV?: Environment) {
        this.parent = parentENV;
        this.variables = new Map();
        this.constants = new Set();
    }

    public declareDefaultVars (env: Environment) {
        env.declareVar("true", MK_BOOL(true), true);
        env.declareVar("false", MK_BOOL(false), true);
        env.declareVar("null", MK_NULL(), true);
    }

    public declareVar (varname: string, value: RuntimeVal, constant = false): RuntimeVal {
        if (this.variables.has(varname)) throw `Cannot declare variable ${varname} as it is already defined within this scope.`;

        this.variables.set(varname, value);
        if (constant) this.constants.add(varname);
        return value;
    }

    public assignVar (varname: string, value: RuntimeVal): RuntimeVal {
        const env = this.resolve(varname);
        if (env.constants.has(varname)) throw `AE: Assignment Error. Cannot reassign value of ${varname} as it is already defined within this scope.`
        env.variables.set(varname, value);
        return value;
    }

    public lookupVar (varname: string): RuntimeVal {
        const env = this.resolve(varname);
        return env.variables.get(varname) as RuntimeVal;
    }

    public resolve (varname: string): Environment {
        if (this.variables.has(varname)) return this;
        if (this.parent == undefined) throw `Cannot resolve '${varname}' within this scope.`
        return this.parent.resolve(varname);
    }

}