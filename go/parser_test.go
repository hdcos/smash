package main

import (
	"reflect"
	"testing"
)

func TestCmdNode(t *testing.T) {
	testTokens := []Token{{column: 0, which: COMMAND, value: "ls"}}

	ast, err := BuildAST(testTokens)
	if err != nil {
		t.Error(err.Error())
	}
	if ast == nil || ast.which != COMMAND {
		t.Error("should build an AST with ls CMD as root")
	}
}

func TestCmdNodeWithArguments(t *testing.T) {
	testTokens := []Token{{column: 0, which: COMMAND, value: "ls"}, {column: 3, which: COMMAND, value: "./folder"}}

	ast, err := BuildAST(testTokens)
	if err != nil {
		t.Error(err.Error())
	}
	if ast == nil || ast.which != COMMAND {
		t.Error("should build an AST with ls CMD as root")
	}
}

func TestAndLeftEmptyCommand(t *testing.T) {
	testTokens := []Token{{column: 0, which: AND, value: "&&"}}

	_, err := BuildAST(testTokens)
	if err == nil {
		t.Error("must raise a syntax error if && without left side ")
	}
}

func TestAndRightEmptyCommand(t *testing.T) {
	testTokens := []Token{
		{column: 0, which: COMMAND, value: "ls"},
		{column: 3, which: AND, value: "&&"},
	}

	_, err := BuildAST(testTokens)
	if err == nil {
		t.Error("must raise a syntax error if && without right side ")
	}
}

func TestAndCommand(t *testing.T) {
	testTokens := []Token{
		{column: 0, which: COMMAND, value: "ls"},
		{column: 0, which: AND, value: "&&"},
		{column: 0, which: COMMAND, value: "wc"},
	}

	ast, err := BuildAST(testTokens)
	if err != nil {
		t.Error(err.Error())
	}
	var expected = &AST{
		which: AND,
		bin:   "",
		args:  []string{},
		children: []*AST{
			{
				which:    COMMAND,
				bin:      "ls",
				args:     []string{},
				children: []*AST{},
			},
			{
				which:    COMMAND,
				bin:      "wc",
				args:     []string{},
				children: []*AST{},
			},
		}}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls && wc\n Got: \n Got: %+v", ast)
	}
}

func TestMultiAndCommand(t *testing.T) {
	testTokens := []Token{
		{column: 0, which: COMMAND, value: "ls"},
		{column: 0, which: AND, value: "&&"},
		{column: 0, which: COMMAND, value: "ls"},
		{column: 0, which: COMMAND, value: "-la"},
		{column: 0, which: AND, value: "&&"},
		{column: 0, which: COMMAND, value: "wc"},
		{column: 0, which: COMMAND, value: "-l"},
	}

	ast, err := BuildAST(testTokens)
	if err != nil {
		t.Error(err.Error())
	}

	var expected = &AST{
		which: AND,
		bin:   "",
		args:  []string{},
		children: []*AST{
			{
				which:    COMMAND,
				bin:      "ls",
				args:     []string{},
				children: []*AST{},
			},
			{
				which:    COMMAND,
				bin:      "ls",
				args:     []string{"-la"},
				children: []*AST{},
			},
			{
				which:    COMMAND,
				bin:      "wc",
				args:     []string{"-l"},
				children: []*AST{},
			},
		}}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls && ls -la && wc -l\n Got: %+v", ast)
	}
}

func TestOrLeftEmptyCommand(t *testing.T) {
	testTokens := []Token{{column: 0, which: OR, value: "||"}}

	_, err := BuildAST(testTokens)
	if err == nil {
		t.Error("must raise a syntax error if && without left side ")
	}
}

func TestOrRightEmptyCommand(t *testing.T) {
	testTokens := []Token{
		{column: 0, which: COMMAND, value: "ls"},
		{column: 3, which: OR, value: "||"},
	}

	_, err := BuildAST(testTokens)
	if err == nil {
		t.Error("must raise a syntax error if && without right side ")
	}
}

func TestMultiOrCommand(t *testing.T) {
	testTokens := []Token{
		{column: 0, which: COMMAND, value: "ls"},
		{column: 0, which: OR, value: "||"},
		{column: 0, which: COMMAND, value: "ls"},
		{column: 0, which: COMMAND, value: "-la"},
		{column: 0, which: OR, value: "||"},
		{column: 0, which: COMMAND, value: "wc"},
		{column: 0, which: COMMAND, value: "-l"},
	}

	ast, err := BuildAST(testTokens)
	if err != nil {
		t.Error(err.Error())
	}

	var expected = &AST{
		which: OR,
		bin:   "",
		args:  []string{},
		children: []*AST{
			{
				which:    COMMAND,
				bin:      "ls",
				args:     []string{},
				children: []*AST{},
			},
			{
				which:    COMMAND,
				bin:      "ls",
				args:     []string{"-la"},
				children: []*AST{},
			},
			{
				which:    COMMAND,
				bin:      "wc",
				args:     []string{"-l"},
				children: []*AST{},
			},
		}}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls || ls -la || wc -l\n Got: %+v", ast)
	}
}

func TestPipeCommand(t *testing.T) {
	testTokens := []Token{
		{column: 0, which: COMMAND, value: "ls"},
		{column: 0, which: PIPE, value: "|"},
		{column: 0, which: COMMAND, value: "wc"},
		{column: 0, which: COMMAND, value: "-l"},
	}

	ast, err := BuildAST(testTokens)
	if err != nil {
		t.Error(err.Error())
	}

	var expected = &AST{
		which: PIPE,
		bin:   "",
		args:  []string{},
		children: []*AST{
			{
				which:    COMMAND,
				bin:      "ls",
				args:     []string{},
				children: []*AST{},
			},
			{
				which:    COMMAND,
				bin:      "wc",
				args:     []string{"-l"},
				children: []*AST{},
			},
		}}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls | wc -l\n Got: [%s %+v]", ast.which, ast.children)
	}
}
