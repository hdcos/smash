package io.smash.parser;

import java.util.List;

import io.smash.lexer.TokenType;

public record AST(
        TokenType which,
        String cmd,
        List<String> args,
        List<AST> children) {
}