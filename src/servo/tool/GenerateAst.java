package servo.tool;

import java.io.IOException;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.List;

public class GenerateAst {

    public static void main(String[] args) throws IOException {
        String outputDir = "/home/wylan/Documents/Projects/ServoLang/Servo/src/servo/ast";

        defineAst(outputDir, "Expr", Arrays.asList(
                "Assign : Token name, Expr value",
                "Binary : Expr left, Token operator, Expr right",
                "Call : Expr callee, Token paren, List<Expr> arguments",
                "Get : Expr object, Token name",
                "Grouping : Expr expression",
                "Literal : Object value",
                "Logical : Expr left, Token operator, Expr right",
                "Set : Expr object, Token name, Expr value",
                "This : Token keyword",
                "Unary : Token operator, Expr right",
                "Variable : Token name"
        ));

        defineAst(outputDir, "Stmt", Arrays.asList(
                "Block : List<Stmt> statements",
                "Class : Token name, List<Stmt.Function> methods",
                "Expression : Expr expression",
                "Function : Token name, List<Token> params, List<Stmt> body, boolean isStatic",
                "If : Expr condition, Stmt thenBranch, Stmt elseBranch",
                "Print : Expr expression",
                "Return : Token keyword, Expr value",
                "Var : Token name, Expr initializer",
                "While : Expr condition, Stmt body"
        ));
    }

    private static void defineAst(String outputDir, String baseName, List<String> types) throws IOException {
        String path = outputDir + "/" + baseName + ".java";
        PrintWriter writer = new PrintWriter(path, StandardCharsets.UTF_8);

        writer.println("package servo.ast;");
        writer.println();
        writer.println("import java.util.List;");
        writer.println();
        writer.println("import servo.token.Token;");
        writer.println();
        writer.println("public abstract class " + baseName + " {");

        defineVisitor(writer, baseName, types);

        for (String type : types) {
            String className = type.split(":")[0].trim();
            String fields = type.split(":")[1].trim();
            defineType(writer, baseName, className, fields);
        }

        writer.println();
        writer.println("    public abstract <R> R accept(Visitor<R> visitor);");

        writer.println("}");
        writer.close();
    }

    private static void defineType(PrintWriter writer, String baseName, String className, String fieldList) {
        writer.println(" public static class " + className + " extends " + baseName + " {");

        // fields
        String[] fields = fieldList.split(", ");
        writer.println();
        for (String field : fields) {
            writer.println("    public final " + field + ";");
        }

        // constructor
        writer.println("    public " + className + "(" + fieldList + ") {");

        //store parameters in fields
        for (String field : fields) {
            String name = field.split(" ")[1];
            writer.println("        this." + name + " = " + name + ";");
        }

        writer.println("    }");

        writer.println();
        writer.println("    @Override");
        writer.println("    public <R> R accept(Visitor<R> visitor) {");
        writer.println("        return visitor.visit" + className + baseName + "(this);");
        writer.println("    }");

        writer.println("    }");
        writer.println();
    }

    private static void defineVisitor(PrintWriter writer, String baseName, List<String> types) {
        writer.println("    public interface Visitor<R> {");

        for (String type : types) {
            String typeName = type.split(":")[0].trim();
            writer.println("        R visit" + typeName + baseName + "(" + typeName + " " + baseName.toLowerCase() + ");");
        }

        writer.println("    }");
    }

}