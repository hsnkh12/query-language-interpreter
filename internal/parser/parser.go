package parser

import (
	"errors"
)

type Parser struct {
	Lexer       Lexer
	Seq         TokenSequence
	nestedDoc   int
	nestedArray int
}

func (p *Parser) Lex() {
	err := p.Lexer.Lex()
	if err != nil {
		panic(err)
	}
}

func (p *Parser) PushToken() {
	p.Seq.Push(p.Lexer.CurrentToken)
}

func (p *Parser) Parse() error {

	p.Lex()
	p.SkipWhiteSpaces()
	var err error

	p.PushToken()

	if p.Lexer.CurrentToken.Type == CREATE {
		err = p.Parse_CreateDelete()
	} else if p.Lexer.CurrentToken.Type == DELETE {
		err = p.Parse_CreateDelete()
	} else if p.Lexer.CurrentToken.Type == RENAME {
		err = p.Parse_Rename()
	} else if p.Lexer.CurrentToken.Type == UPDATE {
	} else if p.Lexer.CurrentToken.Type == ADD {
		err = p.Parse_Add()
	} else if p.Lexer.CurrentToken.Type == GET {
	} else {
		return errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'")
	}

	return err
}

func (p *Parser) Parse_CreateDelete() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.PushToken()

	if p.Lexer.CurrentToken.Type == PROJECT || p.Lexer.CurrentToken.Type == COLLECTION {
		err = p.Parse_CreateDelete_StringL()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'")
	}

	return err
}

func (p *Parser) Parse_Rename() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.PushToken()

	if p.Lexer.CurrentToken.Type == PROJECT || p.Lexer.CurrentToken.Type == COLLECTION {
		err = p.Parse_Rename_StringL()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'")
	}
	return err
}

func (p *Parser) Parse_CreateDelete_StringL() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.End()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

	return err
}

func (p *Parser) Parse_Rename_StringL() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.Parse_Rename_SecStringL()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

	return err
}

func (p *Parser) Parse_Rename_SecStringL() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.End()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expexted string literal")
	}

	return err
}

func (p *Parser) End() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == SEMI_COLUMN {
		err = nil
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ';'")
	}

	return err
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

func (p *Parser) Parse_Add() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == INTO {
		err = p.Parse_Add_StringL()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'into'")
	}

	return err

}

func (p *Parser) Parse_Add_StringL() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.Parse_Add_Doc()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

	return err

}

func (p *Parser) Parse_Add_Doc() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == DOC {
		err = p.Parse_Add_OpenParam()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc'")
	}

	return err
}

func (p *Parser) Parse_Add_OpenParam() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.PushToken()

	if p.Lexer.CurrentToken.Type == OPEN_PARM {
		p.nestedDoc++
		err = p.Parse_Add_Inside_Doc()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc'")
	}

	return err

}

func (p *Parser) Parse_Add_Inside_Doc() error {

	p.Lex()
	var err error
	top := p.Seq.Top().Type

	p.SkipWhiteSpaces()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.PushToken()
			err = p.Parse_Add_Inside_Doc()

		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
		}

	} else if p.Lexer.CurrentToken.Type == NUMBER_LITERAL {

		if top == COLUMN {
			p.PushToken()
			// parse number function
		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore number literal")
		}

	} else if p.Lexer.CurrentToken.Type == OPEN_BRAC {

		if top == COLUMN {
			p.PushToken()
			// parse number function
		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore array")
		}

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.PushToken()
			err = p.Parse_Add_Inside_Doc()
		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair")
		}

	} else if p.Lexer.CurrentToken.Type == COLUMN {

		if top == STRING_LITERAL {
			p.PushToken()
			err = p.Parse_Add_Inside_Doc()
		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key string literal before ':'")
		}

	} else if p.Lexer.CurrentToken.Type == DOC {

		if top == COLUMN {
			p.PushToken()
			err = p.Parse_Add_OpenParam()
		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected ':' berfore 'doc'")
		}

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

		if top == COLUMN || top == COMMA {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value pair before ')'")
		} else {

			p.PushToken()
			if p.nestedDoc > 1 {
				p.nestedDoc--
				err = p.Parse_Add_Inside_Doc()
			} else {
				err = p.End()
			}
		}

	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
	}

	return err
}
