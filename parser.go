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

func NewNode(which string) *AST {
	return &AST{which: which, left: nil, right: nil, bin: "", args: []string{}}
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

func BuildAST(tokens []Token) (*AST, error) {
	expected := []string{COMMAND}

	var root *AST = nil
	i := 0

	for i < len(tokens) {
		currentToken := tokens[i]
		currentTokenType := currentToken.which

		ok := false

		for _, e := range expected {
			if e == currentTokenType {
				ok = true
				break
			}
		}

		if !ok {
			return nil, SyntaxError(currentToken.column, fmt.Sprintf("one of %+v", expected), currentTokenType)
		}

		switch currentTokenType {
		case COMMAND:
			{
				var remainingTokens = tokens[i:]
				commandNode, _ := buildCommandNode(remainingTokens)
				root = commandNode
				expected = []string{AND, OR, PIPE}
				i += 1 + len(commandNode.args) // BIN + ARGS
			}
		case AND, OR, PIPE:
			{
				var remainingTokens = tokens[i+1:]
				if len(remainingTokens) == 0 {
					return nil, SyntaxError(currentToken.column, fmt.Sprintf("%s to have right operand", currentTokenType), "none")
				}
				var newRoot *AST = NewNode(currentTokenType)
				newRoot.left = root
				rightNode, err := buildCommandNode(remainingTokens)
				if err != nil {
					return nil, err
				}
				newRoot.right = rightNode
				root = newRoot
				i += 1 + 1 + len(rightNode.args) // AND/OR/PIPE + BIN + ARGS
				expected = []string{AND, OR, PIPE}
			}
		default:
			{
				i += 1
			}
		}

	}
	return root, nil
}
