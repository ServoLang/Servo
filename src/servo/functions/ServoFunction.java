package servo.functions;

import servo.ast.Stmt;
import servo.classes.ServoInstance;
import servo.env.Environment;
import servo.exception.Return;
import servo.interpreter.Interpreter;

import java.util.List;

public class ServoFunction implements ServoCallable {

    private final Stmt.Function declaration;
    private final Environment closure;

    private final boolean isInitializer;

    public ServoFunction(Stmt.Function declaration, Environment closure, boolean isInitializer) {
        this.isInitializer = isInitializer;
        this.declaration = declaration;
        this.closure = closure;
    }

    public ServoFunction bind(ServoInstance instance) {
        Environment environment = new Environment(getClosure());
        environment.define("this", instance);
        return new ServoFunction(declaration, environment, isInitializer);
    }

    @Override
    public Object call(Interpreter interpreter, List<Object> arguments) {
        Environment environment = new Environment(getClosure());
        for (int i = 0; i < getDeclaration().params.size(); i++) {
            environment.define(getDeclaration().params.get(i).lexeme, arguments.get(i));
        }

        try {
            interpreter.executeBlock(getDeclaration().body, environment);
        } catch (Return returnValue) {
            if (isInitializer) return getClosure().getAt(0, "this");
            return returnValue.value;
        }

        if (isInitializer) return getClosure().getAt(0, "this");
        return null;
    }

    @Override
    public int arity() {
        return getDeclaration().params.size();
    }

    @Override
    public String toString() {
        return "<fn " + getDeclaration().name.lexeme + ">";
    }

    public Stmt.Function getDeclaration() {
        return declaration;
    }

    public Environment getClosure() {
        return closure;
    }

}
