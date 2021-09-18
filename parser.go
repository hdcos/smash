package main

import (
	"fmt"
)

type AST struct {
	left  *AST
	right *AST
	which string
	bin   string
	args  []string
}

func SyntaxError(column int, expected string, found string) error {
	return fmt.Errorf("col[%d]: Expected %s but found %s", column, expected, found)
}

func newNode(which string) *AST {
	return &AST{which: which, left: nil, right: nil, bin: "", args: []string{}}
}

func buildAndNode() *AST {
	return newNode(AND)
}

func buildCommandNode(tokens []Token) (*AST, error) {
	var root *AST = newNode(COMMAND)

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
					return nil, SyntaxError(currentToken.column, "and to have left side", "nothing")
				}
				var remainingTokens = tokens[i+1:]
				if len(remainingTokens) == 0 {
					return nil, SyntaxError(currentToken.column, "and to have right side", "nothing")
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
