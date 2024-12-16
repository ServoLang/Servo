package servo.classes;

import servo.error.RuntimeError;
import servo.functions.ServoFunction;
import servo.token.Token;

import java.util.HashMap;
import java.util.Map;

public class ServoInstance {

    private ServoClass klass;
    private final Map<String, Object> fields = new HashMap<>();

    public ServoInstance(ServoClass klass) {
        this.klass = klass;
    }

    public Object get(Token name) {
        if (fields.containsKey(name.lexeme)) {
            return fields.get(name.lexeme);
        }

        ServoFunction method = klass.findMethod(name.lexeme);
        if (method != null) return method.bind(this);

        throw new RuntimeError(name, "Undefined property '" + name.lexeme + "'.");
    }

    public void set(Token name, Object value) {
        fields.put(name.lexeme, value);
    }

    public ServoClass getKlass() {
        return klass;
    }

    @Override
    public String toString() {
        return klass.getName() + " instance";
    }

}
