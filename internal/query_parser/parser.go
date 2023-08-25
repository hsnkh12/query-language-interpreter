package query_parser

import (
	"errors"
)

type Parser struct {
	Lexer     Lexer
	Seq       TokenSequence
	nestedDoc int
	Err       error
	Stop      bool
}

func CreateNewParser(lexer Lexer) *Parser {
	return &Parser{
		Lexer:     lexer,
		Seq:       TokenSequence{Tokens: []Token{}},
		nestedDoc: 0,
		Err:       nil,
		Stop:      false,
	}
}

func (p *Parser) Lex() {
	p.Lexer.Lex()
	if p.Lexer.Err != nil {
		p.Stop = true
	} else {
		p.SkipWhiteSpaces()
	}
}

func (p *Parser) CurrentTokenType() TokenType {
	return p.Lexer.CurrentToken.Type
}

func (p *Parser) PushToken() {
	p.Seq.Push(p.Lexer.CurrentToken)
}

func (p *Parser) ThrowError(err error) {
	p.Err = err
	p.Stop = true
}

func (p *Parser) Parse() {

	p.Lex()
	p.PushToken()

	switch p.CurrentTokenType() {

	case CREATE:
		p.ParseCreateDelete()
	case DELETE:
		p.ParseCreateDelete()
	case RENAME:
		p.ParseRename()
	case ADD:
		p.ParseAdd()
	case GET:
		p.ParseGet()
	case UPDATE:
		p.ParseUpdate()

	default:
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'"))

	}

}

func (p *Parser) ParseCreateDelete() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	switch p.CurrentTokenType() {

	case PROJECT:
		p.ParseNameStringLiteral()
	case COLLECTION:
		p.ParseNameStringLiteral()
	case FROM:
		if p.Seq.Tokens[0].Type != DELETE {
			break
		}
		p.ParseDeleteDoc()
	default:
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection' keywords"))
		return
	}

	p.ParseEnd()

}

func (p *Parser) ParseRename() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	switch p.CurrentTokenType() {

	case PROJECT:
		p.ParseRenameStringLiteral()
	case COLLECTION:
		p.ParseRenameStringLiteral()
	default:
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection' keywords"))
		return
	}

}

func (p *Parser) ParseRenameStringLiteral() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == STRING_LITERAL {
		p.ParseNameStringLiteral()

		if p.Err == nil {
			p.ParseEnd()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal 'collection_name'"))
	}

}

func (p *Parser) ParseEnd() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() != SEMI_COLUMN {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'"))
		return
	}
}

func (p *Parser) SkipWhiteSpaces() {
	if p.CurrentTokenType() == WHITE_SPACE {
		for {
			p.Lexer.Lex()
			if p.CurrentTokenType() != WHITE_SPACE {
				break
			}
		}
	}
}

func (p *Parser) ParseAdd() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == INTO {
		p.ParseNameStringLiteral()
		if p.Err == nil {
			p.ParseDoc()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'into' keyword"))
		return
	}

}

func (p *Parser) ParseDoc() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == DOC {
		p.ParseOpenParam(true)

		if p.Err == nil {
			p.ParseInsideDoc()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc' keyword"))
		return
	}

}

func (p *Parser) ParseInsideDoc() {

	if p.Stop {
		return
	}

	p.Lex()
	top := p.Seq.Top().Type

	if p.CurrentTokenType() == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.PushToken()
			p.ParseInsideDoc()

		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair"))
			return
		}

	} else if p.CurrentTokenType() == NUMBER_LITERAL || p.CurrentTokenType() == DOT {

		if top == COLUMN {
			p.PushToken()
			p.ParseInsideDoc()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.ParseNumberLiteral()
			p.ParseInsideDoc()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore number literal"))
			return
		}

	} else if p.CurrentTokenType() == COMMA {

		if top != COMMA {
			p.PushToken()
			p.ParseInsideDoc()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair"))
			return
		}

	} else if p.CurrentTokenType() == COLUMN {

		if top == STRING_LITERAL {
			p.PushToken()
			p.ParseInsideDoc()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key string literal before ':'"))
			return
		}

	} else if p.CurrentTokenType() == DOC {

		if top == COLUMN {
			p.PushToken()
			p.ParseOpenParam(true)

			if p.Err == nil {
				p.ParseInsideDoc()
			}
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore 'doc'"))
			return
		}

	} else if p.CurrentTokenType() == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'"))
			return
		} else {

			p.PushToken()
			if p.nestedDoc > 1 {
				p.nestedDoc--
				p.ParseInsideDoc()
			} else {
				p.ParseEnd()
			}
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
		return
	}

}

func (p *Parser) ParseNumberLiteral() {
	p.Seq.ModifyTopLexem(p.Seq.TopLexem() + string(p.Lexer.CurrentLexem))
}

func (p *Parser) ParseGet() {

	if p.Stop {
		return
	}

	p.Lex()

	top := p.Seq.Top().Type
	p.PushToken()

	if p.CurrentTokenType() == ONE {
		if top != ONE {
			p.ParseGet()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'from' keyword"))
			return
		}

	} else if p.CurrentTokenType() == FROM {
		p.ParseNameStringLiteral()

		if p.Err == nil {
			p.ParseAttrs()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'one' or 'from'"))
		return
	}

}

func (p *Parser) ParseNameStringLiteral() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == STRING_LITERAL {
		return
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal"))
		return
	}

}

func (p *Parser) ParseAttrs() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == ATTRS {

		p.ParseOpenParam(false)
		if p.Err == nil {
			p.ParseInsideAttr()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'attrs' keyword"))
		return
	}

}

func (p *Parser) ParseOpenParam(countNested bool) {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == OPEN_PARM {

		if countNested {
			p.nestedDoc++
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected '('"))
		return
	}

}

func (p *Parser) ParseInsideAttr() {

	if p.Stop {
		return
	}

	p.Lex()
	top := p.Seq.Top().Type
	p.PushToken()

	if p.CurrentTokenType() == STRING_LITERAL {

		if top == OPEN_PARM || top == COMMA {
			p.ParseInsideAttr()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ',' before string"))
			return
		}

	} else if p.CurrentTokenType() == COMMA {

		if top != COMMA {
			p.ParseInsideAttr()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string before ','"))
			return
		}
	} else if p.CurrentTokenType() == CLOSE_PARAM {
		p.ParseWhere()
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected literal string attributes"))
		return
	}

}

func (p *Parser) ParseWhere() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == WHERE {

		p.ParseOpenParam(false)
		if p.Err == nil {
			p.ParseInsideWhere()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'attrs' keyword"))
		return
	}

}

func (p *Parser) ParseInsideWhere() {

	if p.Stop {
		return
	}

	p.Lex()
	top := p.Seq.Top().Type

	if p.CurrentTokenType() == STRING_LITERAL {

		if top != STRING_LITERAL && top != NUMBER_LITERAL {
			p.PushToken()
			p.ParseInsideWhere()

		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected logical operation berfore string literal"))
			return
		}

	} else if p.CurrentTokenType() == NUMBER_LITERAL || p.CurrentTokenType() == DOT {

		if IsLogicalOperation(top) {
			p.PushToken()
			p.ParseInsideWhere()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.ParseNumberLiteral()
			p.ParseInsideWhere()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected logical operation berfore number literal"))
			return
		}

	} else if IsLogicalAndOr(p.CurrentTokenType()) {

		if top == STRING_LITERAL || top == NUMBER_LITERAL {
			p.PushToken()
			p.ParseInsideWhere()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', Invalid logical operation"))
			return
		}

	} else if IsLogicalOperation(p.CurrentTokenType()) {

		if !IsLogicalAndOr(top) {
			p.PushToken()
			p.ParseInsideWhere()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', Invalid logical operation"))
			return
		}

	} else if p.CurrentTokenType() == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'"))
			return
		} else {
			p.PushToken()
			p.ParseEnd()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected valid logical operation"))
		return
	}

}

func (p *Parser) ParseUpdate() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == FROM {

		p.ParseNameStringLiteral()

		if p.Err == nil {
			p.ParseSet()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'from' keyword"))
		return
	}
}

func (p *Parser) ParseSet() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == SET {

		p.ParseOpenParam(false)

		if p.Err == nil {
			p.ParseInsideSet()
		}
	}

}

func (p *Parser) ParseInsideSet() {

	if p.Stop {
		return
	}

	p.Lex()
	top := p.Seq.Top().Type

	if p.CurrentTokenType() == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.PushToken()
			p.ParseInsideSet()

		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
			return
		}

	} else if p.CurrentTokenType() == NUMBER_LITERAL || p.CurrentTokenType() == DOT {

		if top == COLUMN {
			p.PushToken()
			p.ParseInsideSet()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.ParseNumberLiteral()
			p.ParseInsideSet()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore number literal"))
			return
		}

	} else if p.CurrentTokenType() == COMMA {

		if top != COMMA {
			p.PushToken()
			p.ParseInsideSet()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair"))
			return
		}

	} else if p.CurrentTokenType() == COLUMN {

		if top == STRING_LITERAL {
			p.PushToken()
			p.ParseInsideSet()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key string literal before ':'"))
			return
		}

	} else if p.CurrentTokenType() == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'"))
			return
		} else {
			p.PushToken()
			p.ParseWhere()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
		return
	}

}

func (p *Parser) ParseDeleteDoc() {

	if p.Stop {
		return
	}

	p.Lex()
	p.PushToken()

	if p.CurrentTokenType() == STRING_LITERAL {
		p.ParseWhere()
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', 'from' keyword"))
		return
	}

}
