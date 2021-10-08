package main

import (
	"reflect"
	"testing"
)

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
	expected := []Token{{which: COMMAND, value: "ls", column: 0}}

	if !reflect.DeepEqual(
		tokens,
		expected) {
		t.Errorf("should tokenize a CMD but got %v", tokens)
	}
}

func TestCommandWithArguments(t *testing.T) {
	typed := "ls -la"
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	expected := []Token{{which: COMMAND, value: "ls", column: 0}, {which: COMMAND, value: "-la", column: 3}}

	if !reflect.DeepEqual(
		tokens,
		expected) {
		t.Errorf("should tokenize a CMD but got %v", tokens)
	}
}

func TestCommandWithPath(t *testing.T) {
	typed := "ls ./F_older#"
	tokens, err := Tokenize(typed)
	if err != nil {
		t.Error(err)
	}
	expected := []Token{{which: COMMAND, value: "ls", column: 0}, {which: COMMAND, value: "./F_older#", column: 3}}

	if !reflect.DeepEqual(
		tokens,
		expected) {
		t.Errorf("should tokenize a CMD but got %v", tokens)
	}
}
