package io.smash.lexer;

public final record Token(TokenType which, int line, int column, CharSequence raw) {
}
