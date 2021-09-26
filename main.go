package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var Debug = false

func debug(prompt string, toDebug interface{}) {
	if Debug {
		fmt.Printf("%s %#v\n", prompt, toDebug)
	}
}

func main() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nBye bye \U0001F44B")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		ok := scanner.Scan()
		if !ok {
			break
		}
		command := scanner.Text()
		if len(command) == 0 {
			continue
		}
		tokens, err := Tokenize(command)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(tokens) == 0 {
			continue
		}
		debug("TOKENS://", tokens)

		ast, err := BuildAST(tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}
		debug("AST://", ast)

		Interpret(ast)
	}

}
