package io.smash.client.cmdline;

import java.util.List;
import java.io.BufferedReader;
import java.io.InputStreamReader;

import io.smash.lexer.Lexer;
import io.smash.lexer.Token;
import io.smash.parser.AST;
import io.smash.parser.Parser;
import io.smash.interpreter.Interpreter;

public class Smash {
    public static void main(String args[]) {
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));

        var exit = false;

        while (!exit) {
            try {
                System.out.print("$> ");
                String input = br.readLine();

                if (input != null) {
                    List<Token> tokens = Lexer.tokenize(input);
                    tokens.forEach((t) -> System.out.println(t));

                    AST parsed = Parser.parse(tokens);
                    if (parsed != null)
                        System.out.println(parsed);

                    Interpreter.exec(parsed);
                } else {
                    exit = true;
                }

            } catch (Exception exception) {
                System.err.println(exception);
            }
        }

        System.out.println("\nbye bye (^ .^)/");
    }
}
