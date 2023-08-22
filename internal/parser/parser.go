package parser

import (
	"errors"
)

type Parser struct {
	Lexer Lexer
	Seq   TokenSequence
}

func (p *Parser) Lex() {
	err := p.Lexer.Lex()
	if err != nil {
		panic(err)
	}
}

func (p *Parser) Parse() error {

	p.Lex()
	p.SkipWhiteSpaces()
	var err error

	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == CREATE {
		err = p.CreateDeleteA()
	} else if p.Lexer.CurrentToken.Type == DELETE {
		err = p.CreateDeleteA()
	} else if p.Lexer.CurrentToken.Type == RENAME {
		err = p.RenameA()
	} else if p.Lexer.CurrentToken.Type == UPDATE {
	} else if p.Lexer.CurrentToken.Type == ADD {
	} else if p.Lexer.CurrentToken.Type == GET {
	} else {
		return errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "'")
	}

	return err
}

func (p *Parser) CreateDeleteA() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == PROJECT || p.Lexer.CurrentToken.Type == COLLECTION {
		err = p.CreateDeleteB()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'")
	}

	return err
}

func (p *Parser) RenameA() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == PROJECT || p.Lexer.CurrentToken.Type == COLLECTION {
		err = p.RenameB()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'project' or 'collection'")
	}
	return err
}

func (p *Parser) CreateDeleteB() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.End()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

	return err
}

func (p *Parser) RenameB() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.RenameC()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

	return err
}

func (p *Parser) RenameC() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()

	p.Seq.Push(p.Lexer.CurrentToken)

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
	p.Seq.Push(p.Lexer.CurrentToken)

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

func (p *Parser) AddA() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == INTO {
		err = p.AddB()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'into'")
	}

	return err

}

func (p *Parser) AddB() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {
		err = p.AddC()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected string literal")
	}

	return err

}

func (p *Parser) AddC() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == DOC {
		err = p.AddD()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc'")
	}

	return err
}

func (p *Parser) AddD() error {

	p.Lex()
	var err error

	p.SkipWhiteSpaces()
	p.Seq.Push(p.Lexer.CurrentToken)

	if p.Lexer.CurrentToken.Type == OPEN_PARM {
		err = p.AddE()
	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected 'doc'")
	}

	return err

}

func (p *Parser) AddE() error {

	p.Lex()
	var err error
	top := p.Seq.Top().Type

	p.SkipWhiteSpaces()

	if p.Lexer.CurrentToken.Type == STRING_LITERAL {

		if top == COMMA || top == OPEN_PARM || top == COLUMN {
			p.Seq.Push(p.Lexer.CurrentToken)
			p.AddE()

		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
		}

	} else if p.Lexer.CurrentToken.Type == OPEN_BRAC {

	} else if p.Lexer.CurrentToken.Type == CLOSED_BRAC {

	} else if p.Lexer.CurrentToken.Type == COMMA {

		if top != COMMA {
			p.AddE()
		} else {
			err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
		}

	} else if p.Lexer.CurrentToken.Type == COLUMN {

	} else if p.Lexer.CurrentToken.Type == DOC {

	} else if p.Lexer.CurrentToken.Type == CLOSE_PARAM {

	} else {
		err = errors.New("PARSER ERROR: Unexpected token: '" + p.Lexer.CurrentToken.Lexem + "', expected key value attributes")
	}

	return err
}
