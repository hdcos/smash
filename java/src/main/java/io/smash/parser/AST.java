package io.smash.parser;

import java.util.List;

import io.smash.lexer.TokenType;

public record AST(
                TokenType which,
                List<String> command,
                List<AST> children) {
}