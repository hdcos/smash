package main

import "fmt"

type Token struct {
	which  string
	value  string
	column int
}

func (t Token) GoString() string {
	return fmt.Sprintf("(@%d %s :: %s)", t.column, t.which, t.value)
}
