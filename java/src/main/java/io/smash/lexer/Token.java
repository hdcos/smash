package io.smash.lexer;

public class Token {
    private String which = null;
    private int lineNumber = 0;
    private int column = 0;

    public int getColumn() {
        return column;
    }

    public int getLineNumber() {
        return lineNumber;
    }

    public String getWhich() {
        return which;
    }

    public Token(String which, int lineNumber, int column) {
        this.which = which;
        this.lineNumber = lineNumber;
        this.column = column;
    }

    public String toString() {
        return String.format("[%d:%d=%s]", this.lineNumber, this.column, this.which);
    }
}
