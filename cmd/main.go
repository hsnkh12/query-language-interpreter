package main

import (
	"fmt"
	"jsondb/internal/parser"
)

func main() {

	src := `add into 'collection' doc('key':doc('key':12,'key':doc('key': doc('key':'value'), 'key':'value'),'key':'value'));`

	lexer_, err := parser.CreateNewLexer(src)

	if err != nil {
		panic(err)
	}

	parser := parser.Parser{Lexer: *lexer_}

	err = parser.Parse()

	if err != nil {
		panic(err)
	}

	for _, tok := range parser.Seq.Tokens {
		fmt.Println(tok)
	}

}
