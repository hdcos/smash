function assertType(value, expected)
    assert(type(value) == expected, 
        string.format("expecting a %s as parameter but got %s", 
            expected,
            type(value))
        )
end
