require("../tokenize")
require("../parse")

describe("Parse", function()
    local cmdLsNode = {id= "CMD", bin = "ls", args = {}, children = {}}

    it("create a simple CMD node without args", function () 
        local expected = cmdLsNode
        local received = Parse(Tokenize("ls"))
        assert.are.same(expected, received)
    end)

    it("create a simple CMD node with args", function () 
        local expected = {id= "CMD", bin = "ls", args = {"-la"}, children = {}}
        local received = Parse(Tokenize("ls -la"))
        assert.are.same(expected, received)
    end)

    it("create an AND ast", function () 
        local expected = {
            id = "AND",
            children = {
                cmdLsNode,
                cmdLsNode,
                cmdLsNode
            }
        }
        local received = Parse(Tokenize("ls && ls && ls"))
        assert.are.same(expected, received)
    end)

    it("create an OR ast", function () 
        local expected = {
            id = "OR",
            children = {
                cmdLsNode,
                cmdLsNode,
                cmdLsNode
            }
        }
        local received = Parse(Tokenize("ls || ls || ls"))
        assert.are.same(expected, received)
    end)

    it("create a PIPE ast", function () 
        local expected = {
            id = "PIPE",
            children = {
                cmdLsNode,
                {id= "CMD", bin = "wc", args = {"-l"}, children = {}}
            }
        }
        local received = Parse(Tokenize("ls | wc -l"))
        assert.are.same(expected, received)
    end)

    it("create a complexe ast", function () 
        local expected = {
            id = "AND",
            children = {
                {
                    id = "PIPE",
                    children = {
                        cmdLsNode,
                        {id= "CMD", bin = "wc", args = {"-l"}, children = {}}
                    }
                },
                cmdLsNode,
            }
        }
        local received = Parse(Tokenize("ls | wc -l && ls"))
        assert.are.same(expected, received)
    end)

end)