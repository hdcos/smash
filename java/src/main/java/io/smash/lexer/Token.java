package io.smash.lexer;

public final record Token(TokenType which, int lineNumber, int column, CharSequence raw) {
}
