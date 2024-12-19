package main

import "errors"

type Parser struct {
	Tokens  []Token
	Current int
}

type Literal struct {
	Kind  TokenType
	Value string
}

func (parser *Parser) peek() Token {
	return parser.Tokens[parser.Current]
}

func (parser *Parser) advance() Token {
	if !parser.atTheEnd() {
		parser.Current++
	}
	token, _ := parser.previous()
	return token
}

func (parser *Parser) atTheEnd() bool {
	return parser.Current+1 == len(parser.Tokens)
}

func (parser *Parser) previous() (Token, error) {
	if parser.Current == 0 {
		return NewToken("", "", ""), errors.New("no previous token is available")
	}
	return parser.Tokens[parser.Current-1], nil
}

func (parser *Parser) peekNext() (Token, error) {
	if !parser.atTheEnd() {
		return parser.Tokens[parser.Current+1], nil
	}

	return NewToken("", "", ""), errors.New("no next token is available")
}

func (parser *Parser) parse() Literal {
	switch parser.peek().Kind {
	case "TRUE":
		expr := Literal{TRUE, "true"}
		return expr
	case "FALSE":
		expr := Literal{FALSE, "false"}
		return expr
	case "NIL":
		return Literal{NIL, "nil"}
	default:
		return Literal{}
	}
}
