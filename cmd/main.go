package main

import (
	"fmt"
	"jsondb/internal/parser"
)

func main() {

	src := `create project 's' ;`

	lexer_, err := parser.CreateNewLexer(src)

	if err != nil {
		panic(err)
	}

	parser := parser.Parser{Lexer: *lexer_}

	err = parser.Parse()

	if err != nil {
		panic(err)
	}

	fmt.Println(parser.Seq)
}
