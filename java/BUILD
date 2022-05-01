load("@rules_java//java:defs.bzl", "java_binary")

java_binary(
    name = "smash",
    srcs = ["src/main/java/io/smash/client/cmdline/Smash.java"],
    main_class = "io.smash.client.cmdline.Smash",
    deps = [
        ":lexer",
        ":parser",
    ],
)

java_library(
    name = "lexer",
    srcs = [
        "src/main/java/io/smash/lexer/Lexems.java",
        "src/main/java/io/smash/lexer/Lexer.java",
        "src/main/java/io/smash/lexer/Token.java",
        "src/main/java/io/smash/lexer/TokenType.java",
        "src/main/java/io/smash/lexer/UnknownCharException.java",
    ],
    visibility = [
        "//src/main/java/io/smash/client/cmdline:__pkg__",
    ],
)

java_library(
    name = "parser",
    srcs = [
        "src/main/java/io/smash/parser/AST.java",
        "src/main/java/io/smash/parser/Parser.java",
    ],
    visibility = [
        "//src/main/java/io/smash/client/cmdline:__pkg__",
    ],
    deps = [
        ":lexer",
    ],
)
