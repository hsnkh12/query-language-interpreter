package parser

type TokenType string

const (
	CREATE          = "CREATE"
	DELETE          = "DELETE"
	RENAME          = "RENAME"
	PROJECT         = "PROJECT"
	TO              = "TO"
	COLLECTION      = "COLLECTION"
	SEMI_COLUMN     = "SEMI_COLUMN"
	COLUMN          = "COLUMN"
	COMMA           = "COMMA"
	ADD             = "ADD"
	INTO            = "INTO"
	DOC             = "DOC"
	OPEN_PARM       = "OPEN_PARM"
	CLOSE_PARAM     = "CLOSE_PARAM"
	ATTRS           = "ATTRS"
	QUOTE           = "QUOTE"
	DOUBLE_QOUTE    = "DOUBLE_QOUTE"
	GET             = "GET"
	FROM            = "FROM"
	WHERE           = "WHERE"
	LOGIC           = "LOGIC"
	LOGICAL_OR      = "LOGICAL_OR"
	LOGICAL_AND     = "LOGICAL_AND"
	LOGICAL_EQUAL   = "LOGICAL_EQUAL"
	LOGICAL_BIGGER  = "LOGICAL_BIGGER"
	LOGICAL_SMALLER = "LOGICAL_SMALLER"
	ONE             = "ONE"
	UPDATE          = "UPDATE"
	SET             = "SET"
	WHITE_SPACE     = "WHITE_SPACE"
	UNKNOWN         = "UNKNOWN"
	STRING_LITERAL  = "STRING_LITERAL"
	NUMBER_LITERAL  = "NUMBER_LITERAL"
	OPEN_BRAC       = "OPEN_BRAC"
	CLOSED_BRAC     = "CLOSED_BRAC"
	DOT             = "DOT"
)

type Token struct {
	Type  TokenType
	Lexem string
}

func tokenize(t TokenType, lexem string) *Token {
	return &Token{Type: t, Lexem: lexem}
}

type TokenSequence struct {
	Tokens []Token
}

func (ts *TokenSequence) Push(token *Token) {
	ts.Tokens = append(ts.Tokens, *token)
}

func (ts *TokenSequence) Top() Token {
	return ts.Tokens[len(ts.Tokens)-1]
}
