package query_interpreter

import (
	"errors"
	qp "jsondb/internal/query_parser"
	"strconv"
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
	var err error = nil

	switch i.Seq.GetCurrentToken().Type {

	case qp.CREATE:
		i.InterpretCreate()
	case qp.DELETE:
		i.InterpretDelete()
	case qp.RENAME:
		i.InterpretRename()
	case qp.ADD:
		err = i.InterpretAdd()
	case qp.GET:
		i.InterpretGet()
	case qp.UPDATE:
		err = i.InterpretUpdate()
	}

	i.Query.Err = err
	return err
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
		i.Query.OPT_TYPE = DELETE_DOCUMENT
		i.Seq.Next()
		i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()
		i.Seq.Next()
		i.InterpretWhere()
		return
	case qp.PROJECT:
		i.Query.OPT_TYPE = DELETE_PROJECT
	case qp.COLLECTION:
		i.Query.OPT_TYPE = DELETE_COLLECTION
	}

	i.Seq.Next()
	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

}

func (i *Interpreter) InterpretGet() {

	i.Seq.Next()

	switch i.Seq.GetCurrentToken().Type {

	case qp.FROM:
		i.Query.OPT_TYPE = GET_ALL_DOCUMENTS
	case qp.ONE:
		i.Query.OPT_TYPE = GET_ONE_DOCUMENT
		i.Seq.Next()
	}

	i.Seq.Next()
	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

	i.Seq.Next()
	i.InterpretAttrs()
	i.Seq.Next()
	i.InterpretWhere()

}

func (i *Interpreter) InterpretUpdate() error {

	i.Query.OPT_TYPE = UPDATE_DOCUMENT
	i.Seq.Next()
	i.Seq.Next()

	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

	i.Seq.Next()
	i.Seq.Next()

	set := []string{}
	values := []interface{}{}

	for i.Seq.GetCurrentToken().Type != qp.CLOSE_PARAM {

		i.Seq.Next()
		set = append(set, i.Seq.GetCurrentLexem())
		i.Seq.Next()
		if i.Seq.GetCurrentToken().Type == qp.CLOSE_PARAM {
			values = append(values, nil)
			break
		}
		if i.Seq.GetCurrentToken().Type == qp.COMMA {
			values = append(values, nil)
			continue
		}
		i.Seq.Next()

		if i.Seq.GetCurrentToken().Type == qp.NUMBER_LITERAL {

			if n, err := strconv.Atoi(i.Seq.GetCurrentLexem()); err == nil {
				values = append(values, n)
			} else if n2, err := strconv.ParseFloat(i.Seq.GetCurrentLexem(), 64); err == nil {
				values = append(values, n2)
			} else {
				return errors.New("INTERPRETER ERROR: invalid number '" + i.Seq.GetCurrentLexem() + "'")

			}

		} else if i.Seq.GetCurrentToken().Type == qp.TRUE {
			values = append(values, true)
		} else if i.Seq.GetCurrentToken().Type == qp.FALSE {
			values = append(values, false)
		} else if i.Seq.GetCurrentToken().Type == qp.NULL {
			values = append(values, nil)
		} else {
			values = append(values, i.Seq.GetCurrentLexem())
		}
		i.Seq.Next()
	}

	i.Query.Kwargs["set"] = set
	i.Query.Kwargs["values"] = values
	i.Seq.Next()
	i.InterpretWhere()

	return nil

}

func (i *Interpreter) InterpretAttrs() {

	i.Seq.Next()
	i.Seq.Next()
	attrs := []string{}

	for i.Seq.GetCurrentToken().Type != qp.CLOSE_PARAM {
		attrs = append(attrs, i.Seq.GetCurrentLexem())
		i.Seq.Next()
		if i.Seq.GetCurrentToken().Type == qp.CLOSE_PARAM {
			break
		}
		i.Seq.Next()
	}

	i.Query.Kwargs["attrs"] = attrs

}

func (i *Interpreter) InterpretWhere() {

	i.Seq.Next()
	i.Seq.Next()

	where := []string{}

	for i.Seq.GetCurrentToken().Type != qp.CLOSE_PARAM {
		where = append(where, i.Seq.GetCurrentLexem())
		i.Seq.Next()
	}
	i.Query.Kwargs["where"] = where

}

func (i *Interpreter) InterpretAdd() error {

	i.Query.OPT_TYPE = qp.ADD

	i.Seq.Next()
	i.Seq.Next()

	i.Query.Kwargs["name"] = i.Seq.GetCurrentLexem()

	i.Seq.Next()
	doc, err := i.InterpretAddDoc()

	if err != nil {
		return err
	}

	i.Query.Kwargs["doc"] = doc

	return nil

}

func (i *Interpreter) InterpretAddDoc() (map[string]interface{}, error) {

	doc := make(map[string]interface{})

	i.Seq.Next()
	var key string

	for i.Seq.GetCurrentToken().Type != qp.CLOSE_PARAM {

		i.Seq.Next()

		if i.Seq.GetCurrentToken().Type == qp.CLOSE_PARAM {
			break
		}
		key = i.Seq.GetCurrentLexem()
		i.Seq.Next()

		if i.Seq.GetCurrentToken().Type == qp.COMMA || i.Seq.GetCurrentToken().Type == qp.CLOSE_PARAM {
			doc[key] = nil
			continue
		}

		i.Seq.Next()

		if i.Seq.GetCurrentToken().Type == qp.DOC {
			v, err := i.InterpretAddDoc()

			if err != nil {
				return nil, err
			}
			doc[key] = v

		} else {

			if i.Seq.GetCurrentToken().Type == qp.NUMBER_LITERAL {

				if n, err := strconv.Atoi(i.Seq.GetCurrentLexem()); err == nil {
					doc[key] = n
				} else if n2, err := strconv.ParseFloat(i.Seq.GetCurrentLexem(), 64); err == nil {
					doc[key] = n2
				} else {
					return nil, errors.New("INTERPRETER ERROR: invalid number '" + i.Seq.GetCurrentLexem() + "'")

				}

			} else if i.Seq.GetCurrentToken().Type == qp.TRUE {
				doc[key] = true
			} else if i.Seq.GetCurrentToken().Type == qp.FALSE {
				doc[key] = false
			} else if i.Seq.GetCurrentToken().Type == qp.NULL {
				doc[key] = nil
			} else {
				doc[key] = i.Seq.GetCurrentLexem()
			}
		}
		i.Seq.Next()

	}

	return doc, nil

}
