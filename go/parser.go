package main

import (
	"fmt"
)

type AST struct {
	children []*AST
	which    string
	bin      string
	args     []string
}

func SyntaxError(column int, expected string, found string) error {
	return fmt.Errorf("col[%d]: Expected %s but found %s", column, expected, found)
}

func NewNode(which string) *AST {
	return &AST{which: which, children: []*AST{}, bin: "", args: []string{}}
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
	var root *AST = nil
	i := 0

	for i < len(tokens) {
		currentToken := tokens[i]
		currentTokenType := currentToken.which

		switch currentTokenType {
		case COMMAND:
			{
				var remainingTokens = tokens[i:]
				commandNode, err := buildCommandNode(remainingTokens)
				if commandNode == nil {
					return nil, err
				}

				if root != nil && (root.which == AND || root.which == OR || root.which == PIPE) {
					root.children = append(root.children, commandNode)
				} else {
					root = commandNode
				}
				i += 1 + len(commandNode.args) // BIN + ARGS
			}
		case AND, OR, PIPE:
			{
				logical := NewNode(currentTokenType)
				if i+1 >= len(tokens) {
					return nil, SyntaxError(i+1, "a command", "EOL")
				}
				if root.which != logical.which {
					logical.children = append(logical.children, root)
					root = logical
				}
				i += 1
			}
		default:
			{
				i += 1
			}
		}

	}
	return root, nil
}
