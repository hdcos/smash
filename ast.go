package main

import (
	"fmt"
)

type AST struct {
	left  *AST
	right *AST
	which string
	bin   string
	args  []string
}

func (a AST) GoString() string {
	if a.which == COMMAND {
		return fmt.Sprintf("%s %+v", a.bin, a.args)
	}
	return fmt.Sprintf("(%s (%s) (%s))", a.which, a.left.GoString(), a.right.GoString())
}
