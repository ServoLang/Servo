package servo.functions;

import servo.interpreter.Interpreter;

import java.util.List;

public interface ServoCallable {

    /**
     * The arity refers to how many arguments a function takes.
     * <p>
     * Example: function testFunc() { ... } has an arity of 0.
     * <p>
     * Example: function newFunc(Integer int, Float f, Object o) { ... } has an arity of 3.
     */
    int arity();

    Object call(Interpreter interpreter, List<Object> arguments);

}
