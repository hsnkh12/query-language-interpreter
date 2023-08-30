package main

import (
	"encoding/json"
	"fmt"
	"jsondb/internal/query_interpreter"
	"jsondb/internal/query_parser"
)

func main() {

	src := `add into 'collection_name' doc('attr1': 'value', 'attr2': 'value', 'attr3' : doc('attr4':'d'));`

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
