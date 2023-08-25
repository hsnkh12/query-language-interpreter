package main

import (
	"fmt"
	"jsondb/internal/query_parser"
)

func main() {

	src := `add into 'collection_name' doc('attr1': 'value', 'attr2': 'value', 'attr3': doc( 'attr4' : 'value'));`

	lexer_, err := query_parser.CreateNewLexer(src)

	if err != nil {
		panic(err)
	}

	parser := query_parser.Parser{Lexer: *lexer_}

	parser.Parse()

	if parser.Err != nil {
		panic(parser.Err)
	}

	for _, tok := range parser.Seq.Tokens {
		fmt.Println(tok)
	}

}
