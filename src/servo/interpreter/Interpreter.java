package servo.interpreter;

import servo.Servo;
import servo.ast.Expr;
import servo.ast.Stmt;
import servo.classes.ServoClass;
import servo.classes.ServoInstance;
import servo.env.Environment;
import servo.error.RuntimeError;
import servo.exception.Return;
import servo.functions.ServoCallable;
import servo.functions.ServoFunction;
import servo.token.Token;
import servo.token.TokenType;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Interpreter implements Expr.Visitor<Object>, Stmt.Visitor<Void> {

    public final Environment globals = new Environment();
    private Environment environment = globals;
    private final Map<Expr, Integer> locals = new HashMap<>();

    public Interpreter() {
        globals.define("clock", new ServoCallable() {
            @Override
            public int arity() {
                return 0;
            }

            @Override
            public Object call(Interpreter interpreter, List<Object> arguments) {
                return (float) System.currentTimeMillis() / 1000.0;
            }

            @Override
            public String toString() {
                return "<native function>";
            }

        });
    }

    public void interpret(List<Stmt> statements) {
        try {
            for (Stmt statement : statements) {
                execute(statement);
            }
        } catch (RuntimeError error) {
            Servo.runtimeError(error);
        }
    }

    private String stringify(Object object) {
        if (object == null) return "null";
        if (object instanceof Integer) {
            String text = object.toString();
            if (text.endsWith(".0")) {
                text = text.substring(0, text.length() - 2);
            }
            return text;
        }

        return object.toString();
    }

    @Override
    public Object visitLiteralExpr(Expr.Literal expr) {
        return expr.value;
    }

    @Override
    public Object visitLogicalExpr(Expr.Logical expr) {
        Object left = evaluate(expr.left);

        if (expr.operator.type == TokenType.OR) {
            if (isTruthy(left)) return left;
        } else {
            if (!isTruthy(left)) return left;
        }

        return evaluate(expr.right);
    }

    @Override
    public Object visitSetExpr(Expr.Set expr) {
        Object object = evaluate(expr.object);

        if (!(object instanceof ServoInstance)) {
            throw new RuntimeError(expr.name, "Only instances have fields.");
        }

        Object value = evaluate(expr.value);
        ((ServoInstance) object).set(expr.name, value);
        return value;
    }

    @Override
    public Object visitSuperExpr(Expr.Super expr) {
        int distance = locals.get(expr);

        ServoClass superclass = (ServoClass) environment.getAt(distance, "super");
        ServoInstance object = (ServoInstance) environment.getAt(distance - 1, "this");
        ServoFunction method = superclass.findMethod(expr.method.lexeme);

        if (method == null) {
            throw new RuntimeError(expr.method, "Undefined property '" + expr.method.lexeme + "'.");
        }

        return method.bind(object);
    }

    @Override
    public Object visitThisExpr(Expr.This expr) {
        return lookupVariable(expr.keyword, expr);
    }

    @Override
    public Object visitUnaryExpr(Expr.Unary expr) {
        Object right = evaluate(expr.right);

        switch (expr.operator.type) {
            case MINUS -> {
                checkNumberOperand(expr.operator, right);
                if (right instanceof Float) {
                    return -(float) right;
                } else if (right instanceof Integer) {
                    return -(int) right;
                }
            }
            case BANG -> {
                return !isTruthy(right);
            }
        }

        // Unreachable
        return null;
    }

    @Override
    public Object visitVariableExpr(Expr.Variable expr) {
        return lookupVariable(expr.name, expr);
    }

    private Object lookupVariable(Token name, Expr expr) {
        Integer distance = locals.get(expr);
        if (distance != null) {
            return environment.getAt(distance, name.lexeme);
        } else {
            return globals.get(name);
        }
    }

    private void checkNumberOperand(Token operator, Object operand) {
        if (operand instanceof Float) return;
        if (operand instanceof Integer) return;
        throw new RuntimeError(operator, "Operand must be a number.");
    }

    private void checkNumberOperands(Token operator, Object left, Object right) {
        if (left instanceof Float && right instanceof Float) return;
        if (left instanceof Integer && right instanceof Integer) return;
        if (left instanceof Float && right instanceof Integer) return;
        if (left instanceof Integer && right instanceof Float) return;
        throw new RuntimeError(operator, "Operands must be numbers.");
    }

    @Override
    public Object visitBinaryExpr(Expr.Binary expr) {
        Object left = evaluate(expr.left);
        Object right = evaluate(expr.right);

        switch (expr.operator.type) {
            case GREATER -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left > (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left > (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left > (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left > (float) right;
                    }
                    default -> {
                    }
                }

            }
            case GREATER_EQUAL -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left >= (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left >= (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left >= (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left >= (float) right;
                    }
                    default -> {
                    }
                }

            }
            case LESS -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left < (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left < (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left < (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left < (float) right;
                    }
                    default -> {
                    }
                }
            }
            case LESS_EQUAL -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left <= (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left <= (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left <= (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left <= (float) right;
                    }
                    default -> {
                    }
                }
            }
            case MINUS -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left - (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left - (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left - (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left - (float) right;
                    }
                    default -> {
                    }
                }

            }
            case PLUS -> {
                if (left instanceof Float && right instanceof Float) {
                    return (float) left + (float) right;
                }

                if (left instanceof Integer && right instanceof Integer) {
                    return (int) left + (int) right;
                }

                if (left instanceof Float && right instanceof Integer) {
                    return (float) left + (int) right;
                }

                if (left instanceof Integer && right instanceof Float) {
                    return (int) left + (float) right;
                }

                if (left instanceof String && right instanceof String) {
                    return left + (String) right;
                }

                if (left instanceof String && right instanceof Integer) {
                    return (String) left + (int) right;
                }

                if (left instanceof String && right instanceof Float) {
                    return (String) left + (float) right;
                }

                throw new RuntimeError(expr.operator, "Operands must be two numbers or two strings.");
            }
            case SLASH -> {
                checkNumberOperands(expr.operator, left, right);
                if (isNaN(left, right)) {
                    Servo.error(expr.operator.line, "NaN: Division by zero.");
                }
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left / (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left / (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left / (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left / (float) right;
                    }
                    default -> {
                    }
                }

            }
            case STAR -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left * (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left * (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left * (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left * (float) right;
                    }
                    default -> {
                    }
                }

            }
            case MODULO -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) left % (float) right;
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) left % (int) right;
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) left % (int) right;
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) left % (float) right;
                    }
                    default -> {
                    }
                }

            }

            case POW -> {
                checkNumberOperands(expr.operator, left, right);
                switch (left) {
                    case Float _ when right instanceof Float -> {
                        return (float) Math.pow((float) left, (float) right);
                    }
                    case Integer _ when right instanceof Integer -> {
                        return (int) Math.pow((int) left, (int) right);
                    }
                    case Float _ when right instanceof Integer -> {
                        return (float) Math.pow((float) left, (int) right);
                    }
                    case Integer _ when right instanceof Float -> {
                        return (int) Math.pow((int) left, (float) right);
                    }
                    default -> {
                    }
                }

            }
            case BANG_EQUAL -> {
                return !isEqual(left, right);
            }
            case EQUAL_EQUAL -> {
                return isEqual(left, right);
            }
        }

        // Unreachable
        return null;
    }

    @Override
    public Object visitCallExpr(Expr.Call expr) {
        Object callee = evaluate(expr.callee);

        List<Object> arguments = new ArrayList<>();
        for (Expr argument : expr.arguments) {
            arguments.add(evaluate(argument));
        }

        if(!(callee instanceof ServoCallable) ) {
            throw new RuntimeError(expr.paren, "Can only call functions and classes.");
        }

        ServoCallable function = (ServoCallable) callee;
        if (arguments.size() != function.arity()) {
            throw new RuntimeError(expr.paren, "Expected " + function.arity() + " arguments but got " + arguments.size() + ".");
        }
        return function.call(this, arguments);
    }

    @Override
    public Object visitGetExpr(Expr.Get expr) {
        Object object = evaluate(expr.object);
        if (object instanceof ServoInstance) {
            return ((ServoInstance) object).get(expr.name);
        }

        throw new RuntimeError(expr.name, "Only instances have properties.");
    }

    @Override
    public Object visitGroupingExpr(Expr.Grouping expr) {
        return evaluate(expr.expression);
    }

    private Object evaluate(Expr expr) {
        return expr.accept(this);
    }

    private void execute(Stmt stmt) {
        stmt.accept(this);
    }

    public void resolve(Expr expr, int depth) {
        locals.put(expr, depth);
    }

    public void executeBlock(List<Stmt> statements, Environment environment) {
        Environment previous = this.environment;
        try {
            this.environment = environment;
            for (Stmt statement : statements) {
                execute(statement);
            }
        } finally {
            this.environment = previous;
        }
    }

    @Override
    public Void visitBlockStmt(Stmt.Block stmt) {
        executeBlock(stmt.statements, new Environment(environment));
        return null;
    }

    @Override
    public Void visitClassStmt(Stmt.Class stmt) {
        Object superclass = null;
        if (stmt.superclass != null) {
            superclass = evaluate(stmt.superclass);
            if (!(superclass instanceof ServoClass)) {
                throw new RuntimeError(stmt.superclass.name, "Superclass must be a class.");
            }
        }

        environment.define(stmt.name.lexeme, null);

        if (stmt.superclass != null) {
            environment = new Environment(environment);
            environment.define("super", superclass);
        }

        Map<String, ServoFunction> methods = new HashMap<>();
        for (Stmt.Function method : stmt.methods) {
            ServoFunction function = new ServoFunction(method, environment, method.name.lexeme.equals("init"));
            methods.put(method.name.lexeme, function);
        }

        ServoClass klass = new ServoClass(stmt.name.lexeme, (ServoClass) superclass, methods);
        if (superclass != null) {
            environment = environment.getEnclosing();
        }
        environment.assign(stmt.name, klass);
        return null;
    }

    @Override
    public Void visitExpressionStmt(Stmt.Expression stmt) {
        evaluate(stmt.expression);
        return null;
    }

    @Override
    public Void visitFunctionStmt(Stmt.Function stmt) {
        ServoFunction function = new ServoFunction(stmt, environment, false);
        environment.define(stmt.name.lexeme, function);
        return null;
    }

    @Override
    public Void visitIfStmt(Stmt.If stmt) {
        if (isTruthy(evaluate(stmt.condition))) {
            execute(stmt.thenBranch);
        } else if (stmt.elseBranch != null) {
            execute(stmt.elseBranch);
        }

        return null;
    }

    // TODO: Supposed to print. Not println. Need to implement other methods.
    @Override
    public Void visitPrintStmt(Stmt.Print stmt) {
        Object value = evaluate(stmt.expression);
        System.out.println(stringify(value));
        return null;
    }

    @Override
    public Void visitReturnStmt(Stmt.Return stmt) {
        Object value = null;
        if (stmt.value != null) value = evaluate(stmt.value);
        throw new Return(value);
    }

    @Override
    public Void visitVarStmt(Stmt.Var stmt) {
        Object value = null;
        if (stmt.initializer != null) {
            value = evaluate(stmt.initializer);
        }
        environment.define(stmt.name.lexeme, value);
        return null;
    }

    @Override
    public Void visitWhileStmt(Stmt.While stmt) {
        while (isTruthy(evaluate(stmt.condition))) {
            execute(stmt.body);
        }
        return null;
    }

    @Override
    public Object visitAssignExpr(Expr.Assign expr) {
        Object value = evaluate(expr.value);

        Integer distance = locals.get(expr);
        if (distance != null) {
            environment.assignAt(distance, expr.name, value);
        } else {
            globals.assign(expr.name, value);
        }

        return value;
    }

    private boolean isTruthy(Object object) {
        if (object == null) return false;
        if (object instanceof Boolean) return (boolean) object;
        return true;
    }

    private boolean isEqual(Object a, Object b) {
        if (a == null && b == null) return true;
        if (a == null) return false;
        return a.equals(b);
    }

    // Returns true if either object is zero.
    private boolean isNaN(Object a, Object b) {
        if (a instanceof Integer && b instanceof Integer) {
            return (int) a == 0 && (int) b == 0;
        }

        if (a instanceof Float && b instanceof Float) {
            return (float) a == 0 && (float) b == 0;
        }
        return false;
    }

}
