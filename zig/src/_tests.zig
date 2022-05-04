const std = @import("std");
const lexer = @import("lexer.zig");
const parser = @import("parser.zig");
const Ast = parser.Ast;
const AstNode = parser.AstNode;
const Token = lexer.Token;
const ArenaAllocator = std.heap.ArenaAllocator;
const alloc = std.testing.allocator;
const expectEqualSlices = std.testing.expectEqualSlices;
const expectEqualStrings = std.testing.expectEqualStrings;
const expectEqual = std.testing.expectEqual;

test "Lexer.tokenize: blank stuff" {
    var tokens = try lexer.tokenize(alloc, "          ");
    defer tokens.deinit();
    try expectEqualSlices(Token, tokens.items, &[_]Token{});
}

test "Lexer.tokenize: AND" {
    var tokens = try lexer.tokenize(alloc, "&&");
    defer tokens.deinit();

    try expectEqualSlices(Token, tokens.items, &[_]Token{
        .{
            .kind = .AND,
            .col = 0,
            .value = "&&",
        },
    });
}

test "Lexer.tokenize: OR" {
    var tokens = try lexer.tokenize(alloc, "||");
    defer tokens.deinit();

    try expectEqualSlices(Token, tokens.items, &[_]Token{
        .{
            .kind = .OR,
            .col = 0,
            .value = "||",
        },
    });
}

test "Lexer.tokenize: PIPE" {
    var tokens = try lexer.tokenize(alloc, "|");
    defer tokens.deinit();

    try expectEqualSlices(Token, tokens.items, &[_]Token{
        .{
            .kind = .PIPE,
            .col = 0,
            .value = "|",
        },
    });
}

test "Lexer.tokenize: CMD" {
    var tokens = try lexer.tokenize(alloc, "ls");
    defer tokens.deinit();

    try expectEqualSlices(Token, tokens.items, &[_]Token{
        .{
            .kind = .CMD,
            .col = 0,
            .value = "ls",
        },
    });
}

test "Lexer.tokenize: CMD (complete)" {
    var tokens = try lexer.tokenize(alloc, "ls -la | cd ../src");
    defer tokens.deinit();

    try expectEqual(tokens.items.len, 5);

    try expectEqual(tokens.items[0].col, 0);
    try expectEqual(tokens.items[0].kind, .CMD);
    try expectEqualStrings(tokens.items[0].value, "ls");
    try expectEqual(tokens.items[1].col, 3);
    try expectEqual(tokens.items[1].kind, .CMD);
    try expectEqualStrings(tokens.items[1].value, "-la");
    try expectEqual(tokens.items[2].col, 7);
    try expectEqual(tokens.items[2].kind, .PIPE);
    try expectEqualStrings(tokens.items[2].value, "|");
    try expectEqual(tokens.items[3].col, 9);
    try expectEqual(tokens.items[3].kind, .CMD);
    try expectEqualStrings(tokens.items[3].value, "cd");
    try expectEqual(tokens.items[4].col, 12);
    try expectEqual(tokens.items[4].kind, .CMD);
    try expectEqualStrings(tokens.items[4].value, "../src");
}

test "Parser.parse: `ls -la -file *.zig`" {
    var tokens = [_]Token{
        .{
            .kind = .CMD,
            .value = "ls",
        },
        .{
            .kind = .CMD,
            .value = "-la",
        },
        .{
            .kind = .CMD,
            .value = "-file *.zig",
        },
    };

    const ast = try parser.parse(alloc, &tokens);
    defer ast.arena.deinit();

    try expectEqual(ast.op_count, 0);
    try expectEqual(ast.cmd_count, 1);
    try expectEqual(ast.args_count, 2);

    try expectEqual(ast.root.children.items.len, 0);
    try expectEqual(ast.root.cmd, "ls");
    try expectEqual(ast.root.args.items.len, 2);
}

test "Parser.parse: `ls -la | cd ../src`" {
    var tokens = [_]Token{
        .{
            .kind = .CMD,
            .value = "ls",
        },
        .{
            .kind = .CMD,
            .value = "-la",
        },
        .{
            .kind = .PIPE,
            .value = "|",
        },
        .{
            .kind = .CMD,
            .value = "cd",
        },
        .{
            .kind = .CMD,
            .value = "../src",
        },
    };

    const ast = try parser.parse(alloc, &tokens);
    defer ast.arena.deinit();

    try expectEqual(ast.op_count, 1);
    try expectEqual(ast.cmd_count, 2);
    try expectEqual(ast.args_count, 2);

    try expectEqual(ast.root.kind, .PIPE);
    try expectEqual(ast.root.children.items.len, 2);
    try expectEqual(ast.root.children.items[0].cmd, "ls");
    try expectEqual(ast.root.children.items[0].args.items[0], "-la");
    try expectEqual(ast.root.children.items[1].cmd, "cd");
    try expectEqual(ast.root.children.items[1].args.items[0], "../src");
}
