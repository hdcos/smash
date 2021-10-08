function Parse(tokens) 
    local root = {id = nil}

    for i,t in pairs(tokens) do
        if t.id == "CMD" then
            local cmdNode = {bin = "", args = {}, children = {}}
            cmdNode.id = "CMD"
            cmdNode.bin = t.value
            cmdNode.args = t.args
            if root.id == "AND" or root.id == "OR" or root.id == "PIPE" then
                root.children[#root.children+1] = cmdNode
            else
                root = cmdNode
            end
        elseif t.id == "AND" then
            if root.id == "AND" then
                -- continue baby
            else 
                local andNode = {id = "AND", children = {root}}
                root = andNode
            end
        elseif t.id == "OR" then
            if root.id == "OR" then
                -- continue baby
            else 
                local orNode = {id = "OR", children = {root}}
                root = orNode
            end
        elseif t.id == "PIPE" then
            if root.id == "PIPE" then
                -- continue baby
            else 
                local pipeNode = {id = "PIPE", children = {root}}
                root = pipeNode
            end
        end
    end
    return root
end