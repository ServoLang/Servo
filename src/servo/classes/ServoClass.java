package servo.classes;

import servo.functions.ServoCallable;
import servo.functions.ServoFunction;
import servo.interpreter.Interpreter;

import java.util.List;
import java.util.Map;

public class ServoClass implements ServoCallable {

    private final String name;
    private final Map<String, ServoFunction> methods;

    public ServoClass(String name, Map<String, ServoFunction> methods) {
        this.name = name;
        this.methods = methods;
    }

    public ServoFunction findMethod(String name) {
        if (getMethods().containsKey(name)) {
            return getMethods().get(name);
        }

        return null;
    }

    public String getName() {
        return name;
    }

    public Map<String, ServoFunction> getMethods() {
        return methods;
    }

    @Override
    public Object call(Interpreter interpreter, List<Object> arguments) {
        ServoInstance instance = new ServoInstance(this);
        ServoFunction initializer = findMethod("init");
        if (initializer != null) initializer.bind(instance).call(interpreter, arguments);
        return instance;
    }

    @Override
    public int arity() {
        ServoFunction initializer = findMethod("init");
        if (initializer == null) return 0;
        return initializer.arity();
    }

    @Override
    public String toString() {
        return name;
    }

}
