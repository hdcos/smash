package io.smash.lexer;

public final record Token(String which, int lineNumber, int column) {
}