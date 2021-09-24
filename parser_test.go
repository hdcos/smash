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
	if ast.which != AND || ast.left.bin != "ls" || ast.right.bin != "wc" {
		t.Error("it should build an ast for ls && wc")
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
		left: &AST{
			which: AND,
			left: &AST{
				which: COMMAND,
				left:  nil,
				right: nil,
				bin:   "ls",
				args:  []string{},
			},
			right: &AST{
				which: COMMAND,
				left:  nil,
				right: nil,
				bin:   "ls",
				args:  []string{"-la"},
			},
			bin:  "",
			args: []string{},
		},
		right: &AST{
			which: COMMAND,
			left:  nil,
			right: nil,
			bin:   "wc",
			args:  []string{"-l"},
		},
		bin:  "",
		args: []string{},
	}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls && ls -la && wc -l\n Got: [%s [%s [%s, %s]] [%s]]", ast.which, ast.left.which, ast.left.left.which, ast.left.right.which, ast.right.which)
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
		left: &AST{
			which: OR,
			left: &AST{
				which: COMMAND,
				left:  nil,
				right: nil,
				bin:   "ls",
				args:  []string{},
			},
			right: &AST{
				which: COMMAND,
				left:  nil,
				right: nil,
				bin:   "ls",
				args:  []string{"-la"},
			},
			bin:  "",
			args: []string{},
		},
		right: &AST{
			which: COMMAND,
			left:  nil,
			right: nil,
			bin:   "wc",
			args:  []string{"-l"},
		},
		bin:  "",
		args: []string{},
	}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls || ls -la || wc -l\n Got: [%s [%s [%s, %s]] [%s]]", ast.which, ast.left.which, ast.left.left.which, ast.left.right.which, ast.right.which)
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
		left: &AST{
			which: COMMAND,
			left:  nil,
			right: nil,
			bin:   "ls",
			args:  []string{},
		},
		right: &AST{
			which: COMMAND,
			left:  nil,
			right: nil,
			bin:   "wc",
			args:  []string{"-l"},
		},
		bin:  "",
		args: []string{},
	}

	if !reflect.DeepEqual(
		ast,
		expected,
	) {
		t.Errorf("should build AST for ls | wc -l\n Got: [%s [%s, %s]]", ast.which, ast.left.which, ast.right.which)
	}
}
