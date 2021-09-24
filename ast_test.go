package main

import "testing"

func TestASTFmt(t *testing.T) {
	a := NewNode(COMMAND)
	a.bin = "ls"
	a.args = []string{"-la"}

	sfmt := a.GoString()
	if sfmt != "ls [-la]" {
		t.Error("should format a AST Node into a string")
	}
}

func TestASTReccursiveFmt(t *testing.T) {
	a := NewNode(AND)
	a.left = NewNode(COMMAND)
	a.left.bin = "ls"
	a.right = NewNode(COMMAND)
	a.right.bin = "wc"

	sfmt := a.GoString()
	if sfmt != "(AND (ls []) (wc []))" {
		t.Error(sfmt)
	}
}
