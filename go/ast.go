package main

import (
	"fmt"
)

func (a AST) GoString() string {
	if a.which == COMMAND {
		return fmt.Sprintf("%s %+v", a.bin, a.args)
	}
	return fmt.Sprintf("(%s (%s) (%s))", a.which, a.left.GoString(), a.right.GoString())
}
