require("../tokenize")

describe("Tokenize", function()
    local lsCmdToken = {id= "CMD", value= "ls", args= {}}

    it("simple command", function () 
        local expected = {lsCmdToken}
        local received = Tokenize("ls")
        assert.are.same(expected, received)
    end)

    it("simple command with args", function () 
        local expected = {{id= "CMD", value= "ls", args={"-la"}}}
        local received = Tokenize("ls -la")
        assert.are.same(expected, received)
    end)

    it("PIPE", function () 
        local expected = {{id= "PIPE", value= "|"}}
        local received = Tokenize(" |  ")
        assert.are.same(expected, received)
    end)

    it("throws if AND not terminated", function () 
        assert.has.error(function() Tokenize("&") end)
        assert.has.error(function() Tokenize("ls & ls ") end)
        assert.has.error(function() Tokenize("ls & ") end)
    end)

    it("logical AND", function () 
        local expected = {
            lsCmdToken,
            {id= "AND", value= "&&"},
            lsCmdToken
        }
        local received = Tokenize("ls && ls")
        assert.are.same(expected, received)
    end)

    it("logical OR", function () 
        local expected = {
            lsCmdToken,
            {id= "OR", value= "||"},
            lsCmdToken
        }
        local received = Tokenize("ls || ls")
        assert.are.same(expected, received)
    end)

    it("multiple logical commands", function () 
        local expected = {
            lsCmdToken,
            {id= "OR", value= "||"},
            lsCmdToken,
            {id= "AND", value= "&&"},
            lsCmdToken
        }
        local received = Tokenize("ls || ls && ls")
        assert.are.same(expected, received)
    end)
end)
