const std = @import("std");
const builtin = @import("builtin");
const Builder = std.build.Builder;
const Pkg = std.build.Pkg;

pub fn build(b: *Builder) void {
    const target = b.standardTargetOptions(.{});
    const mode = b.standardReleaseOptions();

    {
        var test_cmd = b.addTest("src/_tests.zig");
        test_cmd.setTarget(target);
        test_cmd.setBuildMode(mode);

        const test_step = b.step("test", "Run tests");
        test_step.dependOn(&test_cmd.step);
    }

    {
        var exe = b.addExecutable("smash", "src/main.zig");
        exe.setBuildMode(mode);
        exe.setTarget(target);
        exe.install();

        const run_cmd = exe.run();
        run_cmd.step.dependOn(b.getInstallStep());

        const run_step = b.step("run", "Run zig shell");
        run_step.dependOn(&run_cmd.step);
    }
}
