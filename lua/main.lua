require("tokenize")
require("parse")

while true do
    io.write("%> ")
    local line = io.read("l")
    if line == "exit" then break end
    local tokens = Tokenize(line)
    local ast = Parse(tokens)
end