package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

type Output struct {
	success bool
	out     string
}

type EvalContext struct {
	outStream  io.Writer
	errStream  io.Writer
	lastOutput *Output
}

func printOutput(cx *EvalContext) {
	toWrite := []byte(cx.lastOutput.out)
	if cx.lastOutput.success { // sub command failed
		cx.outStream.Write(toWrite)
	} else {
		cx.errStream.Write(append(toWrite, '\n'))
	}
}

func traverse(ast *AST, cx *EvalContext, parent *AST) (*EvalContext, error) {

	switch ast.which {
	case AND:
		{

			cx, err := traverse(ast.left, cx, ast)
			printOutput(cx)
			if !cx.lastOutput.success { // sub command failed
				return cx, err
			} else {
				cx, err := traverse(ast.right, cx, ast)
				printOutput(cx)
				return cx, err
			}
		}
	case OR:
		{
			cx, err := traverse(ast.left, cx, ast)
			printOutput(cx)
			if cx.lastOutput.success { // no need to continue since or
				return cx, err
			} else {
				cx, err := traverse(ast.right, cx, ast)
				printOutput(cx)
				return cx, err
			}
		}
	case PIPE:
		{
			cx, err := traverse(ast.left, cx, ast)

			if cx.lastOutput.success {
				cx, err := traverse(ast.right, cx, ast)
				if parent == nil {
					printOutput(cx)
				}
				return cx, err
			} else {
				return cx, err
			}
		}
	case COMMAND:
		{
			cmd := exec.Command(ast.bin, ast.args...)
			if parent != nil && parent.which == PIPE && cx.lastOutput != nil {
				stdin, err := cmd.StdinPipe()
				if err != nil {
					log.Fatal(err)
				}

				go func() {
					defer stdin.Close()
					io.WriteString(stdin, cx.lastOutput.out)
				}()
			} else {
				cmd.Stdin = os.Stdin
			}
			out, err := cmd.Output()
			commandSucceeded := err == nil
			if commandSucceeded {
				cx.lastOutput = &Output{success: true, out: string(out)}
			} else {
				cx.lastOutput = &Output{success: false, out: string(err.Error())}
			}
			if parent == nil {
				printOutput(cx)
			}
			return cx, nil
		}
	}
	return cx, nil
}

func Interpret(ast *AST) error {

	context := &EvalContext{
		outStream:  os.Stdout,
		errStream:  os.Stderr,
		lastOutput: nil,
	}

	_, err := traverse(ast, context, nil)
	return err
}
