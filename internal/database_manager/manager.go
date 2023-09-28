package databse_manager

import (
	"encoding/json"
	"fmt"
	"jsondb/internal/data_manager"
	qi "jsondb/internal/query_interpreter"
	"jsondb/internal/query_parser"
)

type DatabaseManager struct {
	QuerySrc      string
	ProjectName   string
	MAIN_DIR_PATH string
	Query         *qi.Query
}

func CreateNewDatabaseManager(src string, project_name string, MAIN_DIR_PATH string) *DatabaseManager {
	return &DatabaseManager{
		QuerySrc:      src,
		ProjectName:   project_name,
		MAIN_DIR_PATH: MAIN_DIR_PATH,
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

	var err error

	switch dm.Query.OPT_TYPE {
	case qi.CREATE_PROJECT:
		err = data_manager.CreateProject(dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.DELETE_PROJECT:
		err = data_manager.DeleteProject(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.CREATE_COLLECTION:
		err = data_manager.CreateCollection(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.DELETE_COLLECTION:
		err = data_manager.DeleteCollection(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.RENAME_PROJECT:
		err = data_manager.RenameProject(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.RENAME_COLLECTION:
		err = data_manager.RenamCollection(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.CREATE_DOCUMENT:
		err = data_manager.CreateDocument(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.GET_ALL_DOCUMENTS:
		err = data_manager.GetAllDocuments(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.GET_ONE_DOCUMENT:
		err = data_manager.GetOneDocument(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.UPDATE_DOCUMENT:
		err = data_manager.UpdateDocument(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	case qi.DELETE_DOCUMENT:
		err = data_manager.DeleteDocument(dm.ProjectName, dm.Query.Kwargs, dm.MAIN_DIR_PATH)
	}

	return err
}
