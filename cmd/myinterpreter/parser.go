package main

import (
	"fmt"
	"strings"
)

type Literal struct {
	Value string
}

type Group struct {
	Expr string
}

type Unary struct {
	Operator string
	Expr     string
}

type Statement interface {
	toString()
}

func (literal Literal) toString() {
	fmt.Println(literal.Value)
}

func (group Group) toString() {
	fmt.Printf("(group %v)\n", group.Expr)
}

func (unary Unary) toString() {
	fmt.Printf("(%v %v)\n", unary.Operator, unary.Expr)
}

type Parser struct {
	Tokens  []Token
	Current int
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

func (parser *Parser) parse() Statement {
	switch parser.peek().Lexeme {
	case "(":
		return parser.parse_group()
	case "-", "!":
		return parser.parse_unary()
	default:
		return parser.parse_literal()
	}
}

func (parser *Parser) parse_literal() Literal {
	switch parser.peek().Kind {
	case "TRUE":
		expr := Literal{"true"}
		return expr
	case "FALSE":
		expr := Literal{"false"}
		return expr
	case "NIL":
		return Literal{"nil"}
	case "NUMBER":
		return Literal{parser.peek().Value}
	case "STRING":
		return Literal{parser.peek().Value}
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
			if parser.peek().Lexeme == "(" {
				expr = append(expr, parser.parse_group().Expr)
			}
			expr = append(expr, string(parser.parse_literal().Value))
			parser.advance()
		}
		return Group{strings.Join(expr, "")}
	default:
		return Group{}
	}
}

func (parser *Parser) parse_unary() Unary {
	lexeme := parser.peek().Lexeme
	switch lexeme {
	case "-", "!":
		parser.advance()
		return Unary{string(lexeme), parser.parse_literal().Value}
	default:
		return Unary{}
	}
}
