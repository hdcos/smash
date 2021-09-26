package main

import "testing"

func TestTokensFmt(t *testing.T) {
	token := &Token{which: COMMAND, value: "ls -la", column: 1}

	sfmt := token.GoString()
	if sfmt != "(@1 CMD :: ls -la)" {
		t.Errorf("should format a token into a string %s", sfmt)

	}
}
