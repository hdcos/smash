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
	lsNode := NewNode(COMMAND)
	lsNode.bin = "ls"
	wcNode := NewNode(COMMAND)
	wcNode.bin = "wc"

	a := NewNode(AND)
	a.children = append(a.children, lsNode, wcNode)

	sfmt := a.GoString()
	if sfmt != "{AND (ls [])(wc [])}" {
		t.Error(sfmt)
	}
}
