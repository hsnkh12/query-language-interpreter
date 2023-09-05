package databse_manager

import (
	"encoding/json"
	"fmt"
	qi "jsondb/internal/query_interpreter"
	"jsondb/internal/query_parser"
)

type DatabaseManager struct {
	QuerySrc    string
	ProjectName string
	DIR_PATH    string
	Query       *qi.Query
}

func CreateNewDatabaseManager(src string, project_name string, dir_path string) *DatabaseManager {
	return &DatabaseManager{
		QuerySrc:    src,
		ProjectName: project_name,
		DIR_PATH:    dir_path,
	}
}

func (dm *DatabaseManager) ExecuteQuery() error {

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

	interpreter.Query.Project_name = dm.ProjectName

	dm.Query = interpreter.Query

	q, _ := json.Marshal(interpreter.Query)

	fmt.Println(string(q))

	return nil

}

func (dm *DatabaseManager) CallDataManagerMethod() error {

	switch dm.Query.OPT_TYPE {
	case qi.CREATE_PROJECT:
		return nil
	case qi.DELETE_PROJECT:
		return nil
	case qi.CREATE_COLLECTION:
		return nil
	case qi.DELETE_COLLECTION:
		return nil
	case qi.RENAME_PROJECT:
		return nil
	case qi.RENAME_COLLECTION:
		return nil
	case qi.CREATE_DOCUMENT:
		return nil
	case qi.GET_ALL_DOCUMENTS:
		return nil
	case qi.GET_ONE_DOCUMENT:
		return nil
	case qi.UPDATE_DOCUMENT:
		return nil
	case qi.DELETE_DOCUMENT:
		return nil
	}

	return nil
}
