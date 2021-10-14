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

const BUILTIN_CD = "cd"

func IsBuiltinCommand(value string) bool {
	for _, builtin := range []string{BUILTIN_CD} {
		if builtin == value {
			return true
		}
	}
	return false
}

func execBuiltin(ast *AST, cx *EvalContext) (*EvalContext, error) {
	switch ast.bin {
	case BUILTIN_CD:
		{
			dir, err := os.Getwd()
			if err != nil {
				cx.lastOutput = &Output{success: false, out: string(err.Error())}
				return cx, err
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
				cx.lastOutput = &Output{success: false, out: string(err.Error())}
				return cx, err
			}
			os.Setenv("OLDPWD", dir)
			os.Setenv("PWD", destination)
			cx.lastOutput = &Output{success: true, out: ""}
		}
	}
	return cx, nil
}

func printOutput(cx *EvalContext) {
	toWrite := []byte(cx.lastOutput.out)
	if cx.lastOutput.success { // sub command failed
		cx.outStream.Write(toWrite)
	} else {
		cx.errStream.Write(toWrite)
	}
}

func traverse(ast *AST, cx *EvalContext, parent *AST) (*EvalContext, error) {

	switch ast.which {
	case AND:
		{
			var err error = nil
			for _, c := range ast.children {
				cx, err = traverse(c, cx, ast)
				printOutput(cx)
				if !cx.lastOutput.success { // sub command failed
					return cx, err
				}
			}
			return cx, err
		}
	case OR:
		{
			var err error = nil
			for _, c := range ast.children {
				cx, err = traverse(c, cx, ast)
				printOutput(cx)
				if cx.lastOutput.success { // sub command failed
					return cx, err
				}
			}
			return cx, err
		}
	case PIPE:
		{
			var err error = nil
			for _, c := range ast.children {
				cx, err = traverse(c, cx, ast)
				if !cx.lastOutput.success { // sub command failed
					printOutput(cx)
					cx.lastOutput.out = ""
				}
			}
			printOutput(cx)
			return cx, err
		}
	case COMMAND:
		{
			if IsBuiltinCommand(ast.bin) {
				cx, err := execBuiltin(ast, cx)
				if parent == nil {
					printOutput(cx)
				}
				return cx, err
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
			commandFailed := cmd.ProcessState.ExitCode() == 1
			if err != nil {
				if commandFailed {
					e := err.(*exec.ExitError)
					cx.lastOutput = &Output{success: false, out: string(e.Stderr)}
				} else {
					e := err.(*exec.Error)
					cx.lastOutput = &Output{success: false, out: string(e.Error())}
				}
			} else {
				cx.lastOutput = &Output{success: true, out: string(out)}
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
