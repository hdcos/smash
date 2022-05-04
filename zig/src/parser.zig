const std = @import("std");
const lexer = @import("lexer.zig");
const Allocator = std.mem.Allocator;
const ArenaAllocator = std.heap.ArenaAllocator;
const ArrayList = std.ArrayList;
const Token = lexer.Token;
const TokenType = lexer.TokenType;

pub const Ast = struct {
    /// Root node of our tree.
    root: AstNode,
    /// Mostly used for testing purpose, but all allocated memory
    /// in this struct shared their lifetimes inside this arena allocator.
    arena: ArenaAllocator,
    /// Operators counter.
    op_count: i32 = 0,
    /// Commands counter.
    cmd_count: i32 = 0,
    /// Arguments counter.
    args_count: i32 = 0,
};

pub const AstNode = struct {
    children: ArrayList(AstNode),
    args: ArrayList([]const u8),
    cmd: ?[]const u8 = null,
    kind: TokenType,

    fn init(alloc: Allocator, kind: TokenType) AstNode {
        return .{
            .children = ArrayList(AstNode).init(alloc),
            .args = ArrayList([]const u8).init(alloc),
            .kind = kind,
        };
    }
};

pub fn parse(child_alloc: Allocator, tokens: []const Token) !Ast {
    var ast = Ast{
        .root = undefined,
        .arena = ArenaAllocator.init(child_alloc),
    };

    const alloc = ast.arena.allocator();
    var current: ?AstNode = null;

    var i: usize = 0;
    while (i < tokens.len) {
        const token = tokens[i];

        switch (token.kind) {
            .AND, .OR, .PIPE => {
                var operator_node = AstNode.init(alloc, token.kind);

                if (current) |*node| {
                    if (node.kind != operator_node.kind) {
                        try operator_node.children.append(current.?);
                        current = operator_node;
                    }
                }

                ast.op_count += 1;
                i += 1;
            },
            .CMD => {
                var cmd_node = AstNode.init(alloc, token.kind);

                var delimiter: usize = 0;
                for (tokens[i..]) |item| {
                    if (item.kind != .CMD) break;
                    delimiter += 1;
                }

                for (tokens[i .. i + delimiter]) |next_cmd, index| {
                    if (index == 0) {
                        cmd_node.cmd = next_cmd.value;
                        ast.cmd_count += 1;
                    } else {
                        try cmd_node.args.append(next_cmd.value);
                        ast.args_count += 1;
                    }
                }

                if (current) |*node| {
                    try node.children.append(cmd_node);
                } else {
                    current = cmd_node;
                }

                i += (1 + cmd_node.args.items.len);
            },
        }
    }

    if (current) |node| {
        ast.root = node;
    } else {
        std.debug.panic("Shouldn't be an empty AST.", .{});
    }

    return ast;
}
