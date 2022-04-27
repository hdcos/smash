package io.smash.client.cmdline;

import java.util.List;
import java.io.BufferedReader;
import java.io.InputStreamReader;

import io.smash.lexer.Lexer;
import io.smash.lexer.Token;

public class Smash {
    public static void main(String args[]) {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));

        while (true) {
            try {
                System.out.print("$>");
                String input = br.readLine();

                List<Token> tokens = Lexer.tokenize(input);
                tokens.forEach((t) -> System.out.println(t));
            } catch (Exception exception) {
                System.err.println(exception.toString());
            }
        }

    }
}
