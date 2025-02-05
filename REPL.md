# Read-Eval-Print Loop for Lox

```
for {
    tokens := []
    statements := []
    for {
        source := read()
        tokens.append scan(source)
        for {
             statements.append parse(tokens)
        }
        if 
    }
    result := execute(statements)
    print(result)
}