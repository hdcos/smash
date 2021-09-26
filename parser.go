package main

import (
	"fmt"
)

func SyntaxError(column int, expected string, found string) error {
	return fmt.Errorf("col[%d]: Expected %s but found %s", column, expected, found)
}

func NewNode(which string) *AST {
	return &AST{which: which, left: nil, right: nil, bin: "", args: []string{}}
}

func buildAndNode() *AST {
	return NewNode(AND)
}

func buildOrNode() *AST {
	return NewNode(OR)
}

func buildPipeNode() *AST {
	return NewNode(PIPE)
}

func buildCommandNode(tokens []Token) (*AST, error) {
	var root *AST = NewNode(COMMAND)

	if len(tokens) == 0 || tokens[0].which != COMMAND {
		return nil, SyntaxError(tokens[0].column, "to find COMMAND", tokens[0].which)
	}

	for i := 0; i < len(tokens) && tokens[i].which == COMMAND; i++ {
		if i == 0 {
			root.bin = tokens[i].value
		} else {
			root.args = append(root.args, tokens[i].value)
		}
	}
	return root, nil
}

func buildASTRoot(tokens []Token) (*AST, error) {
	var root *AST = nil
	i := 0

	for i < len(tokens) {
		currentToken := tokens[i]
		switch currentToken.which {
		case COMMAND:
			{
				var remainingTokens = tokens[i:]
				commandNode, _ := buildCommandNode(remainingTokens)
				root = commandNode
				i += 1 + len(commandNode.args) // BIN + ARGS
			}
		case AND:
			{
				if root == nil {
					return nil, SyntaxError(currentToken.column, "AND to have left side", "nothing")
				}
				var remainingTokens = tokens[i+1:]
				if len(remainingTokens) == 0 {
					return nil, SyntaxError(currentToken.column, "AND to have right side", "nothing")
				}
				var newRoot *AST = buildAndNode()
				newRoot.left = root
				rightNode, err := buildCommandNode(remainingTokens)
				if err != nil {
					return nil, err
				}
				newRoot.right = rightNode
				root = newRoot
				i += 1 + 1 + len(rightNode.args) // AND + BIN + ARGS
			}
		case OR:
			{
				if root == nil {
					return nil, SyntaxError(currentToken.column, "OR to have left side", "nothing")
				}
				var remainingTokens = tokens[i+1:]
				if len(remainingTokens) == 0 {
					return nil, SyntaxError(currentToken.column, "OR to have right side", "nothing")
				}
				var newRoot *AST = buildOrNode()
				newRoot.left = root
				rightNode, err := buildCommandNode(remainingTokens)
				if err != nil {
					return nil, err
				}
				newRoot.right = rightNode
				root = newRoot
				i += 1 + 1 + len(rightNode.args) // OR + BIN + ARGS
			}
		case PIPE:
			{
				if root == nil {
					return nil, SyntaxError(currentToken.column, "PIPE to have left side", "nothing")
				}
				var remainingTokens = tokens[i+1:]
				if len(remainingTokens) == 0 {
					return nil, SyntaxError(currentToken.column, "PIPE to have right side", "nothing")
				}
				var newRoot *AST = buildPipeNode()
				newRoot.left = root
				rightNode, err := buildCommandNode(remainingTokens)
				if err != nil {
					return nil, err
				}
				newRoot.right = rightNode
				root = newRoot
				i += 1 + 1 + len(rightNode.args) // OR + BIN + ARGS
			}
		default:
			{
				i++
			}

		}
	}

	return root, nil
}

func BuildAST(tokens []Token) (*AST, error) {

	root, err := buildASTRoot(tokens)
	if err != nil {
		return nil, err
	}

	return root, nil
}
