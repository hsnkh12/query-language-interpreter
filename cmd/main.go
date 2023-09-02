package main

import (
	"encoding/json"
	"fmt"
	"jsondb/internal/query_interpreter"
	"jsondb/internal/query_parser"
)

func main() {

	src := `update from 'collection_name' set('attr1', 'attr2': 'new_value') where();`

	lexer, err := query_parser.CreateNewLexer(src)

	if err != nil {
		panic(err)
	}

	parser := query_parser.CreateNewParser(*lexer)

	interpreter := query_interpreter.CreateNewInterpreter(*parser)

	err = interpreter.Interpret()

	if err != nil {
		panic(err)
	}

	q, _ := json.Marshal(interpreter.Query)

	fmt.Println(string(q))

}
