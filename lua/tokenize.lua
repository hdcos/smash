local function raiseLexicalError(expected, found, position)
    error(string.format("col:%d expecting %s but found %s", position, expected, found), 2)
end

local function raiseUnknownCharError(found, position)
    error(string.format("col:%d unknown char %s", position, found), 2)
end

local function makeToken(id)
    return function (value, ...) 
        return { id = id, value = value } 
    end
end

local function isCmdLexem(c)
    return c >= string.byte("a") and c <= string.byte("z") or c >= string.byte("A") and c <= string.byte("Z") or c == string.byte("-") or c == string.byte(".")
end

local function isWhiteSpace(c)
    return c == string.byte("\t") or c == string.byte(" ")
end

function Tokenize(line)
    local lineLen = string.len(line)
    local k = 1
    local j = 1
    local tokens = {}

    local AND_LEXEM = string.byte("&")
    local PIPE_LEXEM = string.byte("|")

    while j <= lineLen do
        local current = string.byte(line, j)
        if current == AND_LEXEM then
            if j + 1 <= lineLen and string.byte(line, j + 1) == AND_LEXEM then 
                tokens[k] = makeToken("AND")("&&")
                j = j + 2
                k = k + 1
            elseif j + 1 <= lineLen then
                raiseLexicalError("&", string.sub(line, j + 1, j + 1), j + 1)
            else
                raiseLexicalError("&", "EOF", j + 1)
            end
        elseif current == PIPE_LEXEM then
            if j + 1 <= lineLen and string.byte(line, j + 1) == PIPE_LEXEM then 
                tokens[k] = makeToken("OR")("||")
                j = j + 2
                k = k + 1
            else
                tokens[k] = makeToken("PIPE")("|")
                j = j + 1
                k = k + 1
            end
        elseif isCmdLexem(current) then
            local ci = 0
            while j + ci <= lineLen and isCmdLexem(string.byte(line, j + ci )) do
                ci = ci + 1
            end
            local merged = string.sub(line, j, j + ci - 1)
            if k > 1 and tokens[k - 1].id == "CMD" then
                local arg = merged
                local lastToken = tokens[k - 1]
                local lenLastTokenArgs = #lastToken.args
                lastToken.args[lenLastTokenArgs + 1] = merged
            else
                local newCmdToken = makeToken("CMD")(merged)
                newCmdToken.args = {}
                tokens[k] = newCmdToken
            end
            j = j + ci
            k = k + 1
        elseif isWhiteSpace(current) then
            j = j + 1
        else
            raiseUnknownCharError(string.sub(line, j, j), j)
        end
    end
    return tokens
end