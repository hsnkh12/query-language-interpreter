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

func (p *Parser) Lex() {
	p.Lexer.Lex()
	if p.Lexer.Err != nil {
		p.Stop = true
	}
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
	p.SkipWhiteSpaces()
	p.PushToken()

	switch p.Lexer.CurrentToken.Type {

	case CREATE:
		p.Parse_CreateDelete()
	case DELETE:
		p.Parse_CreateDelete()
	case RENAME:
		p.Parse_Rename()
	case ADD:
		p.Parse_Add()
	case GET:
		p.Parse_Get()
	case UPDATE:
		p.Parse_Update()

	default:
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'"))

	}

}

func (p *Parser) Parse_CreateDelete() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	switch p.Lexer.CurrentToken.Type {

	case PROJECT:
		p.Parse_Name_StringLiteral()
	case COLLECTION:
		p.Parse_Name_StringLiteral()
	default:
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'"))
		return
	}

	p.Parse_End()

}

func (p *Parser) Parse_Rename() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	switch p.Lexer.CurrentToken.Type {

	case PROJECT:
		p.Parse_Rename_StringLiteral()
	case COLLECTION:
		p.Parse_Rename_StringLiteral()
	default:
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'"))
		return
	}

}

func (p *Parser) Parse_Rename_StringLiteral() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		p.Parse_Name_StringLiteral()

		if p.Err == nil {
			p.Parse_End()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal"))
	}

}

func (p *Parser) Parse_End() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == SEMI_COLUMN {
		return
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'"))
		return
	}
}

func (p *Parser) SkipWhiteSpaces() {
	if p.Lexer.CurrentToken.Type == WHITE_SPACE {
		for {
			p.Lex()
			if p.Lexer.CurrentToken.Type != WHITE_SPACE {
				break
			}
		}
	}
}

func (p *Parser) Parse_Add() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == INTO {
		p.Parse_Name_StringLiteral()
		if p.Err == nil {
			p.Parse_Doc()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'into'"))
		return
	}

}

func (p *Parser) Parse_Doc() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == DOC {
		p.Parse_OpenParam(true)

		if p.Err == nil {
			p.Parse_Inside_Doc()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc'"))
		return
	}

}

func (p *Parser) Parse_Inside_Doc() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	top := p.Seq.Top().Type

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.PushToken()
			p.Parse_Inside_Doc()

		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == NUMBER_LITERAL || p.Lexer.CurrentToken.Type == DOT {

		if top == COLUMN {
			p.PushToken()
			p.Parse_Inside_Doc()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.Parse_NumberLiteral()
			p.Parse_Inside_Doc()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore number literal"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.PushToken()
			p.Parse_Inside_Doc()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == COLUMN {

		if top == STRING_LITERAL {
			p.PushToken()
			p.Parse_Inside_Doc()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key string literal before ':'"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == DOC {

		if top == COLUMN {
			p.PushToken()
			p.Parse_OpenParam(true)

			if p.Err == nil {
				p.Parse_Inside_Doc()
			}
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore 'doc'"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'"))
			return
		} else {

			p.PushToken()
			if p.nestedDoc > 1 {
				p.nestedDoc--
				p.Parse_Inside_Doc()
			} else {
				p.Parse_End()
			}
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
		return
	}

}

func (p *Parser) Parse_NumberLiteral() {
	p.Seq.ModifyTopLexem(p.Seq.TopLexem() + string(p.Lexer.CurrentLexem))
}

func (p *Parser) Parse_Get() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()

	top := p.Seq.Top().Type
	p.PushToken()

	if p.Lexer.CurrentToken.Type == ONE {
		if top != ONE {
			p.Parse_Get()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'from' keyword"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == FROM {
		p.Parse_Name_StringLiteral()

		if p.Err == nil {
			p.Parse_Attrs()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'one' or 'from'"))
		return
	}

}

func (p *Parser) Parse_Name_StringLiteral() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		return
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal"))
		return
	}

}

func (p *Parser) Parse_Attrs() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == ATTRS {

		p.Parse_OpenParam(false)
		if p.Err == nil {
			p.Parse_Inside_Attr()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'attrs'"))
		return
	}

}

func (p *Parser) Parse_OpenParam(countNested bool) {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == OPEN_PARM {

		if countNested {
			p.nestedDoc++
		}

		return
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected '('"))
		return
	}

}

func (p *Parser) Parse_Inside_Attr() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	top := p.Seq.Top().Type
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == OPEN_PARM || top == COMMA {
			p.Parse_Inside_Attr()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ',' before string"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.Parse_Inside_Attr()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string before ','"))
			return
		}
	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {
		p.Parse_Where()
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected literal string attributes"))
		return
	}

}

func (p *Parser) Parse_Where() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == WHERE {

		p.Parse_OpenParam(false)
		if p.Err == nil {
			p.Parse_Inside_Where()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'attrs'"))
		return
	}

}

func (p *Parser) Parse_Inside_Where() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()

	top := p.Seq.Top().Type

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top != STRING_LITERAL && top != NUMBER_LITERAL {
			p.PushToken()
			p.Parse_Inside_Where()

		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected logical operation berfore string literal"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == NUMBER_LITERAL || p.Lexer.CurrentToken.Type == DOT {

		if IsLogicalOperation(top) {
			p.PushToken()
			p.Parse_Inside_Where()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.Parse_NumberLiteral()
			p.Parse_Inside_Where()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected logical operation berfore number literal"))
			return
		}

	} else if IsLogicalAndOr(p.Lexer.CurrentToken.Type) {

		if top == STRING_LITERAL || top == NUMBER_LITERAL {
			p.PushToken()
			p.Parse_Inside_Where()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', Invalid logical operation"))
			return
		}

	} else if IsLogicalOperation(p.Lexer.CurrentToken.Type) {

		if !IsLogicalAndOr(top) {
			p.PushToken()
			p.Parse_Inside_Where()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', Invalid logical operation"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'"))
			return
		} else {
			p.PushToken()
			p.Parse_End()
		}
	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected valid logical operation"))
		return
	}

}

func (p *Parser) Parse_Update() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == FROM {

		p.Parse_Name_StringLiteral()

		if p.Err == nil {
			p.Parse_Set()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected valid logical operation"))
		return
	}
}

func (p *Parser) Parse_Set() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == SET {

		p.Parse_OpenParam(false)

		if p.Err == nil {
			p.Parse_Inside_Set()
		}
	}

}

func (p *Parser) Parse_Inside_Set() {

	if p.Stop {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()

	top := p.Seq.Top().Type

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.PushToken()
			p.Parse_Inside_Set()

		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == NUMBER_LITERAL || p.Lexer.CurrentToken.Type == DOT {

		if top == COLUMN {
			p.PushToken()
			p.Parse_Inside_Set()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.Parse_NumberLiteral()
			p.Parse_Inside_Set()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore number literal"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.PushToken()
			p.Parse_Inside_Set()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == COLUMN {

		if top == STRING_LITERAL {
			p.PushToken()
			p.Parse_Inside_Set()
		} else {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key string literal before ':'"))
			return
		}

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'"))
			return
		} else {

			p.PushToken()
			p.Parse_Where()
		}

	} else {
		p.ThrowError(errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes"))
		return
	}

}
