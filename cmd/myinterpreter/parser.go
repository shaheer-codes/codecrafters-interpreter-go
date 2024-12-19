package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	Tokens  []Token
	Current int
}

type Literal struct {
	Kind  TokenType
	Value string
}

type Group struct {
	Expr string
}

type Statement interface {
	toString()
}

func (lit Literal) toString() {
	fmt.Println(lit.Value)
}

func (group Group) toString() {
	fmt.Println(group.Expr)
}

func (parser *Parser) peek() Token {
	return parser.Tokens[parser.Current]
}

func (parser *Parser) advance() Token {
	if !parser.atTheEnd() {
		parser.Current++
	}
	token := parser.previous()
	return token
}

func (parser *Parser) atTheEnd() bool {
	return parser.Current+1 == len(parser.Tokens)
}

func (parser *Parser) previous() Token {
	return parser.Tokens[parser.Current-1]
}

func (parser *Parser) peekNext() Token {
	return parser.Tokens[parser.Current+1]
}

func (parser *Parser) parse() Statement {
	switch parser.peek().Kind {
	case "LEFT_PAREN":
		return parser.parse_group()
	default:
		return parser.parse_literal()
	}
}

func (parser *Parser) parse_literal() Literal {
	switch parser.peek().Kind {
	case "TRUE":
		expr := Literal{TRUE, "true"}
		return expr
	case "FALSE":
		expr := Literal{FALSE, "false"}
		return expr
	case "NIL":
		return Literal{NIL, "nil"}
	case "NUMBER":
		return Literal{"NUMBER", parser.peek().Value}
	case "STRING":
		return Literal{"STRING", parser.peek().Value}
	default:
		return Literal{}
	}
}

func (parser *Parser) parse_group() Group {
	switch parser.peek().Lexeme {
	case "(":
		parser.advance()
		var expr []string
		for parser.peek().Lexeme != ")" && parser.peek().Kind != "EOF" {
			expr = append(expr, string(parser.peek().Lexeme))
			parser.advance()
		}
		return Group{fmt.Sprintf("(group %v)", strings.Join(expr, " "))}
	default:
		return Group{}
	}
}
