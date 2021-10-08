package main

import "testing"

func TestIsBuiltin(t *testing.T) {
	if !IsBuiltinCommand(BUILTIN_CD) {
		t.Errorf("cd is a builtin")
	}
	if IsBuiltinCommand("ls") {
		t.Errorf("ls is not a builtin")
	}
}
