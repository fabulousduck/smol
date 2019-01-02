## The lexer

The lexer is the algorithm that takes a program in the form of a string and turns it into tokens. These tokens contain enough data for the AST generator to build a AST out of. 

A token looks like the following

```go
type Token struct {
	Value, Type string
	Line, Col   int
}
```

A main structure for the lexer its self is defined as `Lexer`. This structure is responsible for keeping track of the cursor on the program string. The `Lexer` structure is setup as follows.

```go
type Lexer struct {
	Tokens                                []Token
	currentIndex, currentLine, currentCol int
	FileName, Program                     string
}
```

A lexer can be created with the function `NewLexer(programString string)`.
Example:

```go
    l := NewLexer("MEM A 10")
```

Once you have this lexer structure, you can call the `Lex` function on it.
```go
    l.Lex()
```