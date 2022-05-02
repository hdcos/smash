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

    private static AST buildAndOperator(AST root) {
        var children = new ArrayList<AST>();
        children.add(root);

        return new AST(TokenType.AND, new ArrayList<String>(), children);
    }

    private static AST buildOrOperator(AST root) {
        var children = new ArrayList<AST>();
        children.add(root);

        return new AST(TokenType.OR, new ArrayList<String>(), children);
    }

    private static AST buildPipeOperator(AST root) {
        var children = new ArrayList<AST>();
        children.add(root);

        return new AST(TokenType.PIPE, new ArrayList<String>(), children);
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
                case AND: {
                    if (root.which() != TokenType.AND) {
                        var newRoot = buildAndOperator(root);
                        root = newRoot;
                    } // else just skip this one since it is already chained
                    expected.clear();
                    expected.add(TokenType.CMD);
                    break;
                }

                case OR: {
                    if (root.which() != TokenType.OR) {
                        var newRoot = buildOrOperator(root);
                        root = newRoot;
                    } // else just skip this one since it is already chained
                    expected.clear();
                    expected.add(TokenType.CMD);
                    break;
                }

                case PIPE: {
                    if (root.which() != TokenType.PIPE) {
                        var newRoot = buildPipeOperator(root);
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
