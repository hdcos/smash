const std = @import("std");
const mem = std.mem;
const ascii = std.ascii;
const Allocator = mem.Allocator;
const ArrayList = std.ArrayList;
const TokenIterator = mem.TokenIterator(u8);

const AND = "&&";
const OR = "||";
const PIPE = "|";

pub const TokenType = enum {
    AND,
    OR,
    PIPE,
    CMD,
};

pub const Token = struct {
    kind: TokenType,
    col: usize = 0,
    value: []const u8,
};

const TokenList = ArrayList(Token);

/// Transform given command into a list of tokens.
pub fn tokenize(allocator: Allocator, cmd: []const u8) !TokenList {
    var tokens = TokenList.init(allocator);

    var col: usize = 0;
    var it = mem.split(u8, cmd, " ");

    while (it.next()) |word| {
        if (word.len == 0) continue;

        var token = Token{
            .kind = undefined,
            .col = col,
            .value = word,
        };

        if (mem.eql(u8, word, AND)) {
            token.kind = .AND;
        } else if (mem.eql(u8, word, OR)) {
            token.kind = .OR;
        } else if (mem.eql(u8, word, PIPE)) {
            token.kind = .PIPE;
        } else {
            token.kind = .CMD;
        }

        if (it.index) |index| col = index;

        try tokens.append(token);
    }

    return tokens;
}
