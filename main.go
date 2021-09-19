package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		command := scanner.Text()
		if len(command) == 0 {
			continue
		}
		tokens, err := Tokenize(command)
		if err != nil {
			fmt.Print(err)
			continue
		}

		fmt.Printf("Tokens: %v\n", tokens)

		ast, err := BuildAST(tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("AST: %+v\n", ast)

		Interpret(ast)
	}
}
