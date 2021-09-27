package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

const BUILTIN_CD = "cd"

type Output struct {
	success bool
	out     string
}

type EvalContext struct {
	outStream  io.Writer
	errStream  io.Writer
	lastOutput *Output
}

func IsBuiltinCommand(value string) bool {
	for _, builtin := range []string{BUILTIN_CD} {
		if builtin == value {
			return true
		}
	}
	return false
}

func execBuiltin(ast *AST) (*EvalContext, error) {
	switch ast.bin {
	case BUILTIN_CD:
		{
			dir, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			destination := os.Getenv("HOME")
			if len(ast.args) > 0 {
				destination = ast.args[0]
				if destination == "-" {
					destination = os.Getenv("OLDPWD")
				}
			}
			err = os.Chdir(destination)
			if err != nil {
				return nil, err
			}
			os.Setenv("OLDPWD", dir)
			os.Setenv("PWD", destination)
		}
	}
	return nil, nil
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
			if IsBuiltinCommand(ast.bin) {
				return execBuiltin(ast)
			}
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
