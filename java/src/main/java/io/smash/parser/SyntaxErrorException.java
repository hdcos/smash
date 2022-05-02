package io.smash.parser;

import java.util.ArrayList;
import java.util.List;

import io.smash.lexer.Token;
import io.smash.lexer.TokenType;

public class SyntaxErrorException extends Exception {
    private List<TokenType> expected;
    private Token got;

    public SyntaxErrorException(ArrayList<TokenType> expected, Token got) {
        this.expected = expected;
        this.got = got;
    }

    @Override
    public String toString() {
        return String.format(
                "line:%d column:%d expected one of (%s) but got :: %s",
                this.got.line(),
                this.got.column(),
                this.expected.toString(),
                this.got);
    }
}
