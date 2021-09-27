package main

import (
	"errors"
	"fmt"
)

const AND_SYMBOL = '&'
const OR_SYMBOL = '|'

const AND = "AND"
const OR = "OR"
const COMMAND = "CMD"
const PIPE = "PIPE"

type Token struct {
	which  string
	value  string
	column int
}

func buildLexError(position int, expected byte, found byte) error {
	if found == 0 {
		return errors.New(fmt.Errorf("expecting char %c at %d but found EOF", expected, position).Error())
	}
	return errors.New(fmt.Errorf("expecting char %c at %d but found %c", expected, position, found).Error())
}

func isBlank(b byte) bool {
	return b == ' ' || b == '\t'
}

func isAnd(b byte) bool {
	return b == AND_SYMBOL
}

func isOr(b byte) bool {
	return b == OR_SYMBOL
}

func buildAndToken(s string, i int) (Token, error) {
	j := i + 1
	if j < len(s) {
		if isAnd(s[j]) {
			return Token{which: AND, column: j}, nil
		}
		return Token{}, buildLexError(j, '&', s[j])
	}
	return Token{}, buildLexError(j, '&', 0)
}

func buildOrToken(s string, i int) (Token, error) {
	j := i + 1
	if j < len(s) {
		if isOr(s[j]) {
			return Token{which: OR, column: j}, nil
		}
		return Token{which: PIPE, column: j}, nil
	}
	return Token{}, buildLexError(j, '|', 0)
}

func isCommand(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b == '.' || b == '/'
}

func buildCommandToken(s string, i int) (Token, error) {
	cmd := make([]byte, 0)
	for j := i; j < len(s) && isCommand(s[j]); j++ {
		cmd = append(cmd, s[j])
	}
	return Token{which: COMMAND, value: string(cmd), column: i}, nil
}

func Tokenize(s string) ([]Token, error) {
	res := []Token{}
	if len(s) == 0 {
		return res, nil
	}
	i := 0
	for i < len(s) {
		currentChar := s[i]
		switch {
		case isBlank(currentChar):
			i++
			continue
		case isAnd(currentChar):
			andToken, err := buildAndToken(s, i)
			if err != nil {
				return nil, err
			}
			res = append(res, andToken)
			i += 2
		case isOr(currentChar): // or PIPE
			orToken, err := buildOrToken(s, i)
			if err != nil {
				return nil, err
			}
			res = append(res, orToken)
			var padding = 1 // PIPE
			if orToken.which == OR {
				padding = 2
			}
			i += padding
		case isCommand(currentChar):
			cmdToken, err := buildCommandToken(s, i)
			if err != nil {
				return nil, err
			}
			res = append(res, cmdToken)
			i += len(cmdToken.value)
		default:
			return nil, errors.New(fmt.Errorf("unknown char %c", currentChar).Error())
		}
	}

	return res, nil
}
