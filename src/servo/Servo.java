package servo;

import servo.ast.Stmt;
import servo.error.RuntimeError;
import servo.interpreter.Interpreter;
import servo.parser.Parser;
import servo.scanner.Scanner;
import servo.token.Token;
import servo.token.TokenType;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;

public class Servo {

    static boolean hadError = false;
    static boolean hadRuntimeError = false;

    private static final Interpreter interpreter = new Interpreter();

    public static void main(String[] args) throws IOException {
        if (args.length > 1) {
            System.out.println("Usage: servo [script]");
            System.exit(64);
        } else if (args.length == 1) {
            runFile(args[0]);
        } else {
            //runPrompt();
            runFile("/home/wylan/Documents/Projects/ServoLang/Servo/Test.svo");
        }
    }

    private static void runFile(String filename) throws IOException {
        byte[] bytes = Files.readAllBytes(Paths.get(filename));
        run(new String(bytes, Charset.defaultCharset()));

        if (hadError) System.exit(65);
        if (hadRuntimeError) System.exit(70);
    }

    private static void runPrompt() throws IOException {
        InputStreamReader input = new InputStreamReader(System.in);
        BufferedReader reader = new BufferedReader(input);

        for (;;) {
            System.out.println("> ");
            String line = reader.readLine();
            if (line == null) break;
            run(line);
            hadError = false;
        }
    }

    private static void run (String source) {
        Scanner scanner = new Scanner(source);
        List<Token> tokens = scanner.scanTokens();

        Parser parser = new Parser(tokens);
        List<Stmt> statements = parser.parse();

        // stop if syntax error
        if (hadError) return;

        Resolver resolver = new Resolver(interpreter);
        resolver.resolve(statements);

        // stop if there was a resolution error.
        if (hadError) return;

        interpreter.interpret(statements);
    }

    public static void error(int line, String message) {
        report(line, "", message);
    }

    public static void error(Token token, String message) {
        if (token.type == TokenType.EOF) {
            report(token.line, " at end", message);
        } else {
            report(token.line, " at '" + token.lexeme + "'", message);
        }
    }

    private static void report(int line, String where, String message) {
        System.err.println("[line: " + line + "] Error" + where + ": " + message);
        hadError = true;
    }

    public static void runtimeError(RuntimeError error) {
        System.err.println(error.getMessage() + "\n[line " + error.token.line + ":'" + error.token.lexeme + "']");
        hadRuntimeError = true;
    }

}