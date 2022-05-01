package io.smash.client.cmdline;

import java.util.List;
import java.io.BufferedReader;
import java.io.InputStreamReader;

import io.smash.lexer.Lexer;
import io.smash.lexer.Token;
import io.smash.parser.AST;
import io.smash.parser.Parser;

public class Smash {
    public static void main(String args[]) {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));

        while (true) {
            try {
                System.out.print("$>");
                String input = br.readLine();

                List<Token> tokens = Lexer.tokenize(input);
                tokens.forEach((t) -> System.out.println(t));

                AST parsed = Parser.parse(tokens);
                System.out.println(parsed.toString());

            } catch (Exception exception) {
                System.err.println(exception.toString());
            }
        }

    }
}
