package io.smash.lexer;

public class UnknownCharException extends Exception {
    private int line;
    private int column;
    private char c;

    public UnknownCharException(int line, int column, char c) {
        this.line = line;
        this.column = column;
        this.c = c;
    }

    @Override
    public String toString() {
        return String.format("line:%d column:%d Unknow char '%c'", this.line, this.column, this.c);
    }
}
