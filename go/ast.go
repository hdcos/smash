package main

import (
	"fmt"
)

func (a AST) GoString() string {
	if a.which == COMMAND {
		return fmt.Sprintf("%s %+v", a.bin, a.args)
	}
	res := ""
	for _, c := range a.children {
		res = fmt.Sprintf("%s(%s)", res, c.GoString())
	}

	return fmt.Sprintf("{%s %s}", a.which, res)
}
