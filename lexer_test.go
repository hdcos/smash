package main

import "testing"

func TestBlank(t *testing.T) {
	typed := "    		"
	tokens, _ := Tokenize(typed)
	if len(tokens) != 0 {
		t.Error("empty string should not be tokenized")
	}
}

func TestUnknown(t *testing.T) {
	typed := "ß"
	_, err := Tokenize(typed)
	if err == nil {
		t.Error("should throw when unknown char")
	}
}

func TestAnd(t *testing.T) {
	typed := " && "
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 1 || tokens[0].which != AND {
		t.Error("should tokenize an AND")
	}

	typed = " &"
	_, err = Tokenize(typed)
	if err == nil {
		t.Error("should throw when unmatched AND")
	}
}

func TestOr(t *testing.T) {
	typed := " || "
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 1 || tokens[0].which != OR {
		t.Error("should tokenize an OR")
	}
}

func TestPipe(t *testing.T) {
	typed := " | "
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 1 || tokens[0].which != PIPE {
		t.Error("should tokenize a PIPE")
	}
}

func TestCommand(t *testing.T) {
	typed := "ls"
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 1 || tokens[0].which != COMMAND || tokens[0].value != "ls" {
		t.Error("should tokenize a CMD")
	}
}

func TestCommandWithArguments(t *testing.T) {
	typed := "ls -la"
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 2 || tokens[0].which != COMMAND || tokens[0].value != "ls" || tokens[1].which != COMMAND || tokens[1].value != "-la" {
		t.Error("should tokenize a CMD")
	}
}
