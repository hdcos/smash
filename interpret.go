package main

import (
	"io"
	"os"
	"os/exec"
)

type Output struct {
	success bool
	out     string
}

func traverse(ast *AST, outStream io.Writer) (*Output, error) {

	switch ast.which {
	case AND:
		{

			out, err := traverse(ast.left, outStream)
			if !out.success { // sub command failed
				return out, err
			} else {
				return traverse(ast.right, outStream)
			}
		}
	case COMMAND:
		{
			cmd := exec.Command(ast.bin, ast.args...)
			out, err := cmd.Output()
			commandSucceeded := err == nil
			if commandSucceeded {
				outStream.Write(out)
			}
			return &Output{success: true, out: string(out)}, nil
		}
	}
	return nil, nil
}

func Interpret(ast *AST) error {
	_, err := traverse(ast, os.Stdout)
	return err
}
