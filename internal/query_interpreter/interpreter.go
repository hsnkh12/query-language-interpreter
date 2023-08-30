package query_interpreter

import (
	qp "jsondb/internal/query_parser"
)

type Interpreter struct {
	Parser qp.Parser
	Seq    *qp.TokenSequence
	Query  *Query
}

func CreateNewInterpreter(parser qp.Parser) *Interpreter {
	return &Interpreter{
		Parser: parser,
		Query:  CreateNewQuery(parser.Lexer.Src),
	}
}

func (i *Interpreter) Interpret() error {

	i.Parser.Parse()

	if i.Parser.Lexer.Err != nil {
		i.Query.Err = i.Parser.Lexer.Err
		return i.Query.Err
	}

	if i.Parser.Err != nil {
		i.Query.Err = i.Parser.Err
		return i.Query.Err
	}

	i.Seq = &i.Parser.Seq

	switch i.Seq.GetCurrentToken().Type {

	case qp.CREATE:
		i.InterpretCreate()
	case qp.DELETE:
		i.InterpretDelete()
	case qp.RENAME:
		i.InterpretRename()
	case qp.ADD:
		i.InterpretAdd()
	case qp.GET:
		return nil
	case qp.UPDATE:
		return nil
	}

	return nil
}

func (i *Interpreter) InterpretCreate() {

	i.Seq.Next()

	switch i.Seq.GetCurrentToken().Type {
	case qp.PROJECT:
		i.Query.OPT_TYPE = CREATE_PROJECT
	case qp.COLLECTION:
		i.Query.OPT_TYPE = CREATE_COLLECTION
	}

	i.Seq.Next()
	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

}

func (i *Interpreter) InterpretRename() {

	i.Seq.Next()
	names := make([]string, 2)

	switch i.Seq.GetCurrentToken().Type {
	case qp.PROJECT:
		i.Query.OPT_TYPE = RENAME_PROJECT
	case qp.COLLECTION:
		i.Query.OPT_TYPE = RENAME_COLLECTION
	}

	for x := 0; x < 2; x++ {
		i.Seq.Next()
		names[x] = i.Seq.GetCurrentLexem()
	}

	i.Query.Kwargs["names"] = names
}

func (i *Interpreter) InterpretDelete() {

	i.Seq.Next()

	switch i.Seq.GetCurrentToken().Type {
	case qp.FROM:
		return
	case qp.PROJECT:
		i.Query.OPT_TYPE = DELETE_PROJECT
	case qp.COLLECTION:
		i.Query.OPT_TYPE = DELETE_COLLECTION
	}

	i.Seq.Next()
	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

}

func (i *Interpreter) InterpretAdd() {

	i.Query.OPT_TYPE = qp.ADD

	i.Seq.Next()
	i.Seq.Next()

	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

	doc := make(map[string]interface{})

	i.Seq.Next()
	doc = i.InterpretAddDoc(nil)

	i.Query.Kwargs["doc"] = doc

}

func (i *Interpreter) InterpretAddDoc(doc map[string]interface{}) map[string]interface{} {

	doc = make(map[string]interface{})

	i.Seq.Next()
	var key string

	for i.Seq.GetCurrentToken().Type != qp.CLOSE_PARAM {

		i.Seq.Next()

		key = i.Seq.GetCurrentLexem()
		i.Seq.Next()
		i.Seq.Next()

		if i.Seq.GetCurrentToken().Type == qp.DOC {
			i.InterpretAddDoc(doc)
		}

		doc[key] = i.Seq.GetCurrentLexem()
		i.Seq.Next()

	}

	return doc

}
