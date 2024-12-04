import Parser from "./frontend/parser.ts";
import Environment from "./runtime/environment.ts";
import { evaluate } from "./runtime/interpreter.ts";

repl();

function repl () {
    const parser = new Parser();
    const env = new Environment();
    env.declareDefaultVars(env);

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