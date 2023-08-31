package query_parser

import (
	"errors"
	"strconv"
)

type Lexer struct {
	Src          string
	Position     int
	CurrentChar  rune
	CurrentToken *Token
	CurrentLexem string
	Err          error
}

func CreateNewLexer(src string) (*Lexer, error) {
	if src[len(src)-1] != ';' {
		return nil, errors.New("LEXER ERROR: Missing ';' at the end of the line")
	}
	return &Lexer{
		Src:          src,
		Position:     -1,
		CurrentChar:  rune(src[0]),
		CurrentToken: nil,
		CurrentLexem: "",
	}, nil
}

func (l *Lexer) Lex() {

	src := l.Src

	if l.Position >= len(src)-1 {
		return
	}

	l.Position++
	l.CurrentChar = rune(src[l.Position])
	var tokenType TokenType
	startAt := l.Position

	for {

		if l.CurrentChar == ';' {
			if tokenType == UNKNOWN {
				break
			}
			tokenType = SEMI_COLUMN
			l.CurrentLexem = ";"
			break
		}

		if l.CurrentChar == ' ' {
			if tokenType == UNKNOWN {
				break
			}
			tokenType = WHITE_SPACE
			l.CurrentLexem = " "
			break
		}
		l.CurrentLexem = src[startAt : l.Position+1]

		tokenType = DetectTokenType(l.CurrentLexem)

		if tokenType != UNKNOWN {
			break
		}

		l.Position++

		l.CurrentChar = rune(src[l.Position])
	}

	token := tokenize(tokenType, l.CurrentLexem)
	l.CurrentToken = token

	if l.CurrentToken.Type == UNKNOWN {
		l.Err = errors.New("LEXER ERROR: Unknown token '" + l.CurrentLexem + "'")
		return
	}
}

func DetectTokenType(lexem string) TokenType {

	switch lexem {
	case "create", "CREATE":
		return CREATE
	case "delete", "DELETE":
		return DELETE
	case "rename", "RENAME":
		return RENAME
	case "project", "PROJECT":
		return PROJECT
	case "to", "TO":
		return TO
	case "collection", "COLLECTION":
		return COLLECTION
	case "add", "ADD":
		return ADD
	case "into", "INTO":
		return INTO
	case "get", "GET":
		return GET
	case "from", "FROM":
		return FROM
	case "one", "ONE":
		return ONE
	case "update", "UPDATE":
		return UPDATE
	case "doc", "DOC":
		return DOC
	case "attrs", "ATTRS":
		return ATTRS
	case "where", "WHERE":
		return WHERE
	case "logic", "LOGIC":
		return LOGIC
	case "set", "SET":
		return SET
	case "||":
		return LOGICAL_OR
	case "&&":
		return LOGICAL_AND
	case "==":
		return LOGICAL_EQUAL
	case ">":
		return LOGICAL_BIGGER
	case "<":
		return LOGICAL_SMALLER
	case ";":
		return SEMI_COLUMN
	case ":":
		return COLUMN
	case ",":
		return COMMA
	case "(":
		return OPEN_PARM
	case ")":
		return CLOSE_PARAM
	case "\"":
		return DOUBLE_QOUTE
	case "[":
		return OPEN_BRAC
	case "]":
		return CLOSED_BRAC
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "null":
		return NULL
	case ".":
		return DOT
	case " ":
		return WHITE_SPACE
	default:

		if len(lexem) >= 2 && lexem[0] == '\'' && lexem[len(lexem)-1] == '\'' {
			return STRING_LITERAL
		}
		if _, err := strconv.Atoi(lexem); err == nil {
			return NUMBER_LITERAL
		}
		return UNKNOWN
	}
}
