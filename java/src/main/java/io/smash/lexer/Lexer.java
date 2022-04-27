package io.smash.lexer;

import java.util.List;
import java.util.Vector;

public class Lexer {

    public static boolean isAndLexem(char c) {
        return c == Lexems.AND_CHAR;
    }

    public static boolean isOrLexem(char c) {
        return c == Lexems.OR_CHAR;
    }

    public static boolean isBlank(char c) {
        return c == ' ' || c == '\t';
    }

    public static boolean isAlphaNum(char c) {
        return (c >= 'a' && c <= 'z')
                || (c >= 'A' && c <= 'Z')
                || c == '-'
                || (c >= '0' && c <= '9');
    }

    public static List<Token> tokenize(String line) throws UnknownCharException {
        List<Token> tokens = new Vector<Token>();

        final int lineLength = line.length();

        if (line != null && lineLength > 0) {
            int i = 0;

            while (i < lineLength) {
                char c = line.charAt(i);

                if (isBlank(c)) {
                    i++;
                } else if (isAndLexem(c)) {
                    if (i + 1 < lineLength && isAndLexem(line.charAt(i + 1))) {
                        final Token and = new Token(TokenType.AND, 0, i, Lexems.AND);
                        tokens.add(and);
                        i += 2;
                    } else {
                        throw new UnknownCharException(0, i, c);
                    }
                } else if (isOrLexem(c)) {
                    if (i + 1 < lineLength && isOrLexem(line.charAt(i + 1))) {
                        final Token or = new Token(TokenType.OR, 0, i, Lexems.OR);
                        tokens.add(or);
                        i += 2;
                    } else {
                        final Token pipe = new Token(TokenType.PIPE, 0, i, Lexems.PIPE);
                        tokens.add(pipe);
                        i++;
                    }
                } else if (isAlphaNum(c)) {
                    StringBuilder sb = new StringBuilder();
                    int j = 0;
                    char cc = c;
                    while (isAlphaNum(cc)) {
                        sb.append(cc);
                        j++;
                        if (i + j < lineLength) {
                            cc = line.charAt(i + j);
                        } else
                            break;
                    }
                    final Token cmdPart = new Token(TokenType.CMD, 0, i, sb.subSequence(0, j));
                    tokens.add(cmdPart);
                    i += sb.length();
                } else {
                    throw new UnknownCharException(0, i, c);
                }
            }

        }

        return tokens;
    }
}
