package query_interpreter

import (
	"math/rand"
	"time"
)

type QueryOperationType string

const (
	CREATE_PROJECT    = "CREATE_PROJECT"
	DELETE_PROJECT    = "DELETE_PROJECT"
	CREATE_COLLECTION = "CREATE_COLLECTION"
	DELETE_COLLECTION = "DELETE_COLLECTION"
	RENAME_PROJECT    = "RENAME_PROJECT"
	RENAME_COLLECTION = "RENAME_COLLECTION"
	CREATE_DOCUMENT   = "CREATE_DOCUMENT"
	GET_ALL_DOCUMENTS = "GET_ALL_DOCUMENTS"
	GET_ONE_DOCUMENT  = "GET_ONE_DOCUMENT"
	UPDATE_DOCUMENT   = "UPDATE_DOCUMENT"
	DELETE_DOCUMENT   = "DELETE_DOCUMENT"
)

type Query struct {
	Id             uint32                 `json:"query_id"`
	QueryStatement string                 `json:"query_statement"`
	OPT_TYPE       QueryOperationType     `json:"operation_type"`
	Project_name   string                 `json:"project_name"`
	Kwargs         map[string]interface{} `json:"kwargs"`
	Response       *string                `json:"response"`
	Err            error                  `json:"error"`
}

func CreateNewQuery(src string) *Query {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Uint32()

	return &Query{
		Id:             randomNumber,
		QueryStatement: src,
		Kwargs:         make(map[string]interface{}),
	}
}
