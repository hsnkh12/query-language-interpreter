package querymanager

import (
	"encoding/json"
	"fmt"
	qi "jsondb/internal/query_interpreter"
	"jsondb/internal/query_parser"
)

type QueryManager struct {
	QuerySrc string
	Query    *qi.Query
}

func New(src string) *QueryManager {
	return &QueryManager{
		QuerySrc: src,
	}
}

func (dm *QueryManager) ExecuteQuery() error {

	lexer, err := query_parser.CreateNewLexer(dm.QuerySrc)

	if err != nil {
		return err
	}

	parser := query_parser.CreateNewParser(*lexer)

	interpreter := qi.CreateNewInterpreter(*parser)

	err = interpreter.Interpret()

	if err != nil {
		return err
	}

	dm.Query = interpreter.Query

	q, _ := json.Marshal(interpreter.Query)

	fmt.Println(string(q))

	return nil

}
