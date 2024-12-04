import Parser from "./frontend/parser.ts";
import Environment from "./runtime/environment.ts";
import { evaluate } from "./runtime/interpreter.ts";
import { MK_NULL, MK_BOOL } from "./runtime/values.ts";

repl();

function repl () {
    const parser = new Parser();
    const env = new Environment();
    env.declareVar("true", MK_BOOL(true), true);
    env.declareVar("false", MK_BOOL(false), true);
    env.declareVar("null", MK_NULL(), true);

    console.log("\nServo-Repl v0.0.1");
    while (true) {
        const input = prompt("> ");
        // Exit
        if (!input || input.includes("exit")) {
            Deno.exit(1);
        }

        const program = parser.produceAST(input);
        const result = evaluate(program, env);
        console.log(result);
    }
}