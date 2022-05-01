package io.smash.parser;

import java.util.List;

import io.smash.lexer.Token;

public class Parser {

    public static AST parse(List<Token> tokens) {
        if (tokens == null || tokens.size() == 0) {
            return null;
        }

        Token first = tokens.get(0);
        AST root = new AST(first.which(), first.raw().toString(), null, null);

        return root;
    }
}
