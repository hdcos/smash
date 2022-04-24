package io.smash.lexer;

import java.util.List;
import java.util.Vector;

public class Lexer {

    public static boolean isAndLexem(char c) {
        return c == Lexems.AND_CHAR;
    }

    public static boolean isBlank(char c) {
        return c == ' ' || c == '\t';
    }

    public static List<Token> tokenize(String line) {
        List<Token> tokens = new Vector<Token>();

        final int lineLength = line.length();

        char previous = '\0';

        if (line != null && lineLength > 0) {
            int i = 0;

            while (i < lineLength) {
                char c = line.charAt(i);

                if (!isBlank(c)) {
                    if (isAndLexem(c) && isAndLexem(previous)) {
                        final Token and = new Token(Lexems.AND, 0, i);
                        tokens.add(and);
                    }
                }

                previous = c;
                i++;
            }

        }

        return tokens;
    }
}
