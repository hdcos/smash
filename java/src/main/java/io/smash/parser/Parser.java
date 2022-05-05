package io.smash.parser;

import java.util.ArrayList;
import java.util.List;

import io.smash.lexer.Token;
import io.smash.lexer.TokenType;

public class Parser {

    private static AST buildCommand(AST root, List<Token> cmd) {

        var commandNode = new AST(TokenType.CMD, new ArrayList<String>(), null);

        for (int i = 0; i < cmd.size() && cmd.get(i).which() == TokenType.CMD; i++) {
            commandNode.command().add(cmd.get(i).raw().toString());
        }

        // this is the called command name
        return commandNode;
    }

    private static AST buildOperator(AST root, TokenType type) {
        var children = new ArrayList<AST>();
        children.add(root);

        return new AST(type, new ArrayList<String>(), children);
    }

    public static AST parse(List<Token> tokens) throws SyntaxErrorException {
        if (tokens == null || tokens.size() == 0) {
            return null;
        }

        AST root = null;

        ArrayList<TokenType> expected = new ArrayList<TokenType>();
        expected.add(TokenType.CMD);

        for (var i = 0; i < tokens.size(); i++) {
            Token currentToken = tokens.get(i);
            if (!expected.contains(currentToken.which())) {
                throw new SyntaxErrorException(expected, currentToken);
            }

            switch (currentToken.which()) {
                case AND:
                case PIPE:
                case OR: {
                    if (root.which() != TokenType.AND) {
                        var newRoot = buildOperator(root, currentToken.which());
                        root = newRoot;
                    } // else just skip this one since it is already chained
                    expected.clear();
                    expected.add(TokenType.CMD);
                    break;
                }

                case CMD: {
                    var rest = tokens.subList(i, tokens.size());
                    var commandNode = buildCommand(root, rest);

                    if (root != null
                            && (root.which() == TokenType.AND
                                    || root.which() == TokenType.OR
                                    || root.which() == TokenType.PIPE)) {
                        root.children().add(commandNode);
                    } else {
                        root = commandNode;
                    }

                    expected.clear();
                    expected.add(TokenType.AND);
                    expected.add(TokenType.OR);
                    expected.add(TokenType.PIPE);

                    i += commandNode.command().size() - 1; // -1 because i++ at end of iteration;
                    break;
                }

                default:
                    break;
            }
        }

        return root;
    }
}
