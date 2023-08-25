package query_parser

import (
	"errors"
)

type Parser struct {
	Lexer     Lexer
	Seq       TokenSequence
	nestedDoc int
	Err       error
}

func (p *Parser) Lex() {
	err := p.Lexer.Lex()
	if err != nil {
		p.Err = err
	}
}

func (p *Parser) PushToken() {
	p.Seq.Push(p.Lexer.CurrentToken)
}

func (p *Parser) Parse() {

	p.Lex()
	p.SkipWhiteSpaces()

	p.PushToken()

	if p.Lexer.CurrentToken.Type == CREATE {
		p.Parse_CreateDelete()
	} else if p.Lexer.CurrentToken.Type == DELETE {
		p.Parse_CreateDelete()
	} else if p.Lexer.CurrentToken.Type == RENAME {
		p.Parse_Rename()
	} else if p.Lexer.CurrentToken.Type == UPDATE {
	} else if p.Lexer.CurrentToken.Type == ADD {
		p.Parse_Add()
	} else if p.Lexer.CurrentToken.Type == GET {
		p.Parse_Get()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'")
	}

}

func (p *Parser) Parse_CreateDelete() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == PROJECT || p.Lexer.CurrentToken.Type == COLLECTION {
		p.Parse_CreateDelete_StringL()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'")
		return
	}
}

func (p *Parser) Parse_Rename() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == PROJECT || p.Lexer.CurrentToken.Type == COLLECTION {
		p.Parse_Rename_StringL()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'")
		return
	}
}

func (p *Parser) Parse_CreateDelete_StringL() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		p.Parse_End()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
		return
	}
}

func (p *Parser) Parse_Rename_StringL() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		p.Parse_Rename_SecStringL()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

}

func (p *Parser) Parse_Rename_SecStringL() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		p.Parse_End()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expexted string literal")
		return
	}
}

func (p *Parser) Parse_End() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == SEMI_COLUMN {
		return
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'")
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

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == INTO {
		p.Parse_Name_StringLiteral()
		if p.Err == nil {
			p.Parse_Add_Doc()
		}
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'into'")
		return
	}

}

func (p *Parser) Parse_Add_Doc() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == DOC {
		p.Parse_OpenParam(true)

		if p.Err == nil {
			p.Parse_Add_Inside_Doc()
		}
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc'")
		return
	}

}

func (p *Parser) Parse_Add_Inside_Doc() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	top := p.Seq.Top().Type

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.PushToken()
			p.Parse_Add_Inside_Doc()

		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
			return
		}

	} else if p.Lexer.CurrentToken.Type == NUMBER_LITERAL || p.Lexer.CurrentToken.Type == DOT {

		if top == COLUMN {
			p.PushToken()
			p.Parse_Add_Inside_Doc()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.Parse_NumberLiteral()
			p.Parse_Add_Inside_Doc()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore number literal")
			return
		}

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.PushToken()
			p.Parse_Add_Inside_Doc()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair")
			return
		}

	} else if p.Lexer.CurrentToken.Type == COLUMN {

		if top == STRING_LITERAL {
			p.PushToken()
			p.Parse_Add_Inside_Doc()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key string literal before ':'")
			return
		}

	} else if p.Lexer.CurrentToken.Type == DOC {

		if top == COLUMN {
			p.PushToken()
			p.Parse_OpenParam(true)

			if p.Err == nil {
				p.Parse_Add_Inside_Doc()
			}
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore 'doc'")
			return
		}

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'")
			return
		} else {

			p.PushToken()
			if p.nestedDoc > 1 {
				p.nestedDoc--
				p.Parse_Add_Inside_Doc()
			} else {
				p.Parse_End()
			}
		}

	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
		return
	}

}

func (p *Parser) Parse_NumberLiteral() {
	p.Seq.ModifyTopLexem(p.Seq.TopLexem() + string(p.Lexer.CurrentLexem))
}

func (p *Parser) Parse_Get() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()

	top := p.Seq.Top().Type

	if p.Lexer.CurrentToken.Type == ONE {
		if top != ONE {
			p.Parse_Get()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected from")
			return
		}

	} else if p.Lexer.CurrentToken.Type == FROM {
		p.Parse_Name_StringLiteral()

		if p.Err == nil {
			p.Parse_Get_Attrs()
		}

	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'one' or 'from'")
		return
	}

}

func (p *Parser) Parse_Name_StringLiteral() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		return
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
		return
	}

}

func (p *Parser) Parse_Get_Attrs() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == ATTRS {

		p.Parse_OpenParam(false)
		if p.Err == nil {
			p.Parse_Get_Inside_Attr()
		}

	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'attrs'")
		return
	}

}

func (p *Parser) Parse_OpenParam(isDoc bool) {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == OPEN_PARM {

		if isDoc {
			p.nestedDoc++
		}

		return
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected '('")
		return
	}

}

func (p *Parser) Parse_Get_Inside_Attr() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	top := p.Seq.Top().Type

	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == OPEN_PARM || top == COMMA {
			p.Parse_Get_Inside_Attr()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ',' before string")
			return
		}

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.Parse_Get_Inside_Attr()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string before ','")
			return
		}
	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {
		p.Parse_Get_Where()
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected literal string attributes")
		return
	}

}

func (p *Parser) Parse_Get_Where() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == WHERE {

		p.Parse_OpenParam(false)
		if p.Err == nil {
			p.Parse_Add_Inside_Doc()
		}

	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'attrs'")
		return
	}

}

func (p *Parser) Parse_Get_Inside_Where() {

	if p.Err != nil {
		return
	}

	p.Lex()
	p.SkipWhiteSpaces()

	top := p.Seq.Top().Type

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top != STRING_LITERAL && top != NUMBER_LITERAL {
			p.PushToken()
			p.Parse_Get_Inside_Where()

		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
			return
		}

	} else if p.Lexer.CurrentToken.Type == NUMBER_LITERAL || p.Lexer.CurrentToken.Type == DOT {

		if IsLogicalOperation(top) {
			p.PushToken()
			p.Parse_Get_Inside_Where()

		} else if top == NUMBER_LITERAL || top == DOT {
			p.Parse_NumberLiteral()
			p.Parse_Get_Inside_Where()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected logical operation berfore number literal")
			return
		}

	} else if IsLogicalAndOr(p.Lexer.CurrentToken.Type) {

		if top == STRING_LITERAL || top == NUMBER_LITERAL {
			p.PushToken()
			p.Parse_Get_Inside_Where()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', Invalid logical operation")
			return
		}

	} else if IsLogicalOperation(p.Lexer.CurrentToken.Type) {

		if !IsLogicalAndOr(top) {
			p.PushToken()
			p.Parse_Get_Inside_Where()
		} else {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', Invalid logical operation")
			return
		}

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'")
			return
		} else {

			p.Parse_End()
		}
	} else {
		p.Err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected logical operation")
		return
	}

}
