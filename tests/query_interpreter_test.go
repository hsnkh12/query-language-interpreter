package tests

import (
	"encoding/json"
	"jsondb/internal/query_interpreter"
	"jsondb/internal/query_parser"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		query  string
		hasErr bool
	}{
		// Test cases for CREATE
		{"create project 'project_name';", false},
		{"create collection 'collection_name';", false},

		// Test cases for DELETE
		{"delete project 'project_name';", false},
		{"delete collection 'collection_name';", false},

		// Test cases for RENAME
		{"rename project 'project_name' 'new_name';", false},
		{"rename collection 'collection_name' 'new_name';", false},

		// Test cases for ADD
		{"add into 'collection_name' doc('attr1': 'value', 'attr2': 'value', 'attr3': doc( 'attr4' : 'value'));", false},

		// Test cases for GET
		{"get from 'collection_name' attrs('attr1', 'attr2') where('attr1' == 'attr2' || 'attr2' > 'attr4');", false},
		{"get one from 'collection_name' attrs() where();", false},

		// Test cases for UPDATE
		{"update from 'collection_name' set('attr1': 'new_value', 'attr2': 'new_value') where();", false},

		// Test cases for DELETE
		{"delete from 'collection_name' where();", false},
	}

	for _, test := range tests {
		lexer, _ := query_parser.CreateNewLexer(test.query)

		parser := query_parser.CreateNewParser(*lexer)

		interpreter := query_interpreter.CreateNewInterpreter(*parser)

		err := interpreter.Interpret()

		if err != nil {
			t.Errorf("Query: %s\nExpected no error, but got: %v", test.query, err)
		} else if lexer.Err != nil {
			t.Errorf("Query: %s\nExpected no error, but got: %v", test.query, lexer.Err)
		} else if parser.Err != nil && !test.hasErr {
			t.Errorf("Query: %s\nExpected no error, but got: %v", test.query, parser.Err)
		} else if parser.Err == nil && test.hasErr {
			t.Errorf("Query: %s\nExpected error, but got none", test.query)
		}

		j, _ := json.Marshal(interpreter.Query)
		t.Logf("%s OUTPUT: %+v", test.query, string(j))
		t.Log("\n----------------------------\n")
	}
}
