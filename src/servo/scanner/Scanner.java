package servo.scanner;

import servo.Servo;
import servo.token.Token;
import servo.token.TokenType;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static servo.token.TokenType.*;

public class Scanner {

    private static final Map<String, TokenType> keywords;

    private final String source;
    private final List<Token> tokens = new ArrayList<>();
    private int start = 0;
    private int current = 0;
    private int line = 1;

    public Scanner(String source) {
        this.source = source;
    }

    public List<Token> scanTokens() {
        while (!isAtEnd()) {
            start = current;
            scanToken();
        }

        tokens.add(new Token(EOF, "", null, line));
        return tokens;
    }

    private boolean isAtEnd() {
        return current >= source.length();
    }

    private void scanToken() {
        char c = advance();
        switch (c) {
            case '(': {
                addToken(LEFT_PAREN);
            }
            break;

            case ')': {
                addToken(RIGHT_PAREN);
            }
            break;

            case '{': {
                addToken(LEFT_BRACE);
            }
            break;

            case '}': {
                addToken(RIGHT_BRACE);
            }
            break;

            case ',': {
                addToken(COMMA);
            }
            break;

            case '.': {
                addToken(DOT);
            }
            break;

            case '-': {
                addToken(match('>') ? RIGHT_POINTER : MINUS);
            }
            break;

            case '+': {
                addToken(PLUS);
            }
            break;

            case ';': {
                addToken(SEMICOLON);
            }
            break;

            case '*': {
                addToken(STAR);
            }
            break;

            case '^': {
                addToken(POW);
            }
            break;

            case '%': {
                addToken(MODULO);
            }
            break;

            case '!': {
                addToken(match('=') ? BANG_EQUAL : BANG);
            }
            break;

            case '=': {
                addToken(match('=') ? EQUAL_EQUAL : EQUAL);
            }
            break;

            case '<': {
                if (match('=')) {
                    addToken(LESS_EQUAL);
                } else if (match('-')) {
                    addToken(LEFT_POINTER);
                } else {
                    addToken(LESS);

                }
            }
            break;

            case '>': {
                addToken(match('=') ? GREATER_EQUAL : GREATER);
            }
            break;

            case '/': {
                if (match('/')) {
                    while (peek() != '\n' && !isAtEnd()) advance();
                } else if (match('*')) { //multi line comment operator
                    while (peek() != '*' && peekNext() != '/' && !isAtEnd()) advance();
                } else {
                    addToken(SLASH);
                }
            }
            break;

            case ' ':
            case '\r':
            case '\t':
                break;

            case '\n': {
                line++;
            }
            break;

            case '`':
            case '\'':
            case '"': {
                string(c);
            }
            break;

            default:
                if (isDigit(c)) {
                    number();
                } else if (isAlpha(c)) {
                    identifier();
                } else {
                    Servo.error(line, "Unexpected token at: " + line);
                }
                break;
        }
    }

    private boolean match(char expected) {
        if (isAtEnd()) return false;
        if (source.charAt(current) != expected) return false;
        current++;
        return true;
    }

    private char peek() {
        if (isAtEnd()) return '\0';
        return source.charAt(current);
    }

    // TODO: Allow leading decimals
    private boolean isDigit(char c) {
        return c >= '0' && c <= '9';
    }

    private boolean isAlpha(char c) {
        return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_';
    }

    private boolean isAlphaNumeric(char c) {
        return isAlpha(c) || isDigit(c);
    }

    private void identifier() {
        while (isAlphaNumeric(peek())) advance();
        String text = source.substring(start, current);
        TokenType type = keywords.get(text);
        if (type == null) type = IDENTIFIER;
        addToken(type);
    }

    // Differentiates integers and floats
    private void number() {
        while (isDigit(peek())) advance();

        if (peek() == '.' && isDigit(peekNext())) {
            do advance();
            while (isDigit(peek()));
            addToken(FLOAT, Float.parseFloat(source.substring(start, current)));
        } else {
            addToken(INTEGER, Integer.parseInt(source.substring(start, current)));
        }
    }

    private char peekNext() {
        if (current + 1 >= source.length()) return '\0';
        return source.charAt(current + 1);
    }

    private char advance() {
        return source.charAt(current++);
    }

    private void addToken(TokenType type) {
        addToken(type, null);
    }

    private void addToken(TokenType type, Object literal) {
        String text = source.substring(start, current);
        tokens.add(new Token(type, text, literal, line));
    }

    private void string(char end) {
        while (peek() != end && !isAtEnd()) {
            if (peek() == '\n') line++;
            advance();
        }

        if (isAtEnd()) {
            Servo.error(line, "Unterminated string at: " + line);
            return;
        }

        advance();
        String value = source.substring(start + 1, current - 1);
        addToken(STRING, value);
    }

    static {
        keywords = new HashMap<>();
        keywords.put("and", AND);
        keywords.put("class", CLASS);
        keywords.put("static", STATIC);
        keywords.put("else", ELSE);
        keywords.put("false", FALSE);
        keywords.put("for", FOR);
        keywords.put("function", FUNCTION);
        keywords.put("if", IF);
        keywords.put("null", NULL);
        keywords.put("or", OR);
        keywords.put("print", PRINT);
        keywords.put("printf", PRINTF);
        keywords.put("printl", PRINTL);
        keywords.put("printlf", PRINTLF);
        keywords.put("return", RETURN);
        keywords.put("super", SUPER);
        keywords.put("this", THIS);
        keywords.put("true", TRUE);
        keywords.put("var", VAR);
        keywords.put("let", LET);
        keywords.put("const", CONST);
        keywords.put("constant", CONSTANT);
        keywords.put("while", WHILE);
        keywords.put("scope", SCOPE);
        keywords.put("export", EXPORT);
        keywords.put("interface", INTERFACE);
        keywords.put("enum", ENUM);
        keywords.put("struct", STRUCT);
        keywords.put("allocate", ALLOCATE);
        keywords.put("deallocate", DEALLOCATE);
        keywords.put("reference", REFERENCE);
        keywords.put("Integer", INTEGER);
        keywords.put("Integer64", INTEGER64);
        keywords.put("Float", FLOAT);
        keywords.put("Float64", FLOAT64);
        keywords.put("String", STRING);
        keywords.put("Boolean", BOOLEAN);
        keywords.put("Byte", BYTE);
        keywords.put("Char", CHAR);
    }

}
