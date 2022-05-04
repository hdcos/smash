const std = @import("std");
const lexer = @import("lexer.zig");
const parser = @import("parser.zig");
const ArrayList = std.ArrayList;
const ArenaAllocator = std.heap.ArenaAllocator;
const io = std.io;
const os = std.os;
const mem = std.mem;
const PATH = [std.fs.MAX_PATH_BYTES]u8;

const stdin = io.getStdIn().reader();
const stdout = io.getStdOut().writer();

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    var is_running = true;

    while (is_running) {
        // Wrap all allocation during a single frame.
        var arena = ArenaAllocator.init(gpa.allocator());
        const alloc = arena.allocator();

        // Note: Now, we don't have to think about memory at all.
        defer arena.deinit();

        var buf: PATH = undefined;
        try stdout.print("{s}> ", .{try os.getcwd(&buf)});

        if (try stdin.readUntilDelimiterOrEofAlloc(alloc, '\n', 1000)) |cmd| {
            if (mem.eql(u8, cmd, "exit")) {
                is_running = false;
                continue;
            }

            const tokens = try lexer.tokenize(alloc, cmd);
            _ = try parser.parse(alloc, tokens.items);
        }
    }

    try stdout.print("Bye!\n", .{});
}
