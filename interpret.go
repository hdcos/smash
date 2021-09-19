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

func traverse(ast *AST, outStream io.Writer, errStream io.Writer) (*Output, error) {

	switch ast.which {
	case AND:
		{

			out, err := traverse(ast.left, outStream, errStream)
			if !out.success { // sub command failed
				return out, err
			} else {
				return traverse(ast.right, outStream, errStream)
			}
		}
	case OR:
		{
			out, err := traverse(ast.left, outStream, errStream)
			if out.success { // sub command failed
				return out, err
			} else {
				return traverse(ast.right, outStream, errStream)
			}
		}
	case COMMAND:
		{
			cmd := exec.Command(ast.bin, ast.args...)
			out, err := cmd.Output()
			commandSucceeded := err == nil
			if commandSucceeded {
				outStream.Write(out)
				return &Output{success: true, out: string(out)}, nil
			} else {
				errStream.Write(append([]byte(err.Error()), '\n'))
				return &Output{success: false, out: string(out)}, nil
			}
		}
	}
	return nil, nil
}

func Interpret(ast *AST) error {
	_, err := traverse(ast, os.Stdout, os.Stderr)
	return err
}
