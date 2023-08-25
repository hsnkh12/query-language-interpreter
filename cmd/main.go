package main

import (
	"fmt"
	"jsondb/internal/query_parser"
)

func main() {

	src := `update from 'collection_name' set('attr1': 'new_value', 'attr2': 'new_value') where('attr' == 'd');`

	lexer_, err := query_parser.CreateNewLexer(src)

	if err != nil {
		panic(err)
	}

	parser := query_parser.Parser{Lexer: *lexer_}

	parser.Parse()

	if lexer_.Err != nil {
		panic(lexer_.Err)
	}

	if parser.Err != nil {
		panic(parser.Err)
	}

	for _, tok := range parser.Seq.Tokens {
		fmt.Println(tok)
	}

}
