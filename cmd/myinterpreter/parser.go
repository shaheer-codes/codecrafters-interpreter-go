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

type Binary struct {
	Left     string
	Operator string
	Right    string
}

type Statement interface {
	toString() string
}

func (literal Literal) toString() string {
	return literal.Value
}

func (group Group) toString() string {
	return fmt.Sprintf("(group %v)", group.Expr)
}

func (unary Unary) toString() string {
	return fmt.Sprintf("(%v %v)", unary.Operator, unary.Expr)
}

func (binary Binary) toString() string {
	return fmt.Sprintf("(%v %v %v)", binary.Operator, binary.Left, binary.Right)
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
	case "*", "/":
		return parser.parse_binary()
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
				expr = append(expr, parser.parse_group().toString())
			} else if parser.peek().Lexeme == "-" || parser.peek().Lexeme == "!" {
				expr = append(expr, parser.parse_unary().toString())
			}
			expr = append(expr, string(parser.parse_literal().toString()))
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
		nextToken := parser.peek()
		if nextToken.Value == "null" {
			if nextToken.Lexeme == "(" {
				return Unary{string(lexeme), parser.parse_group().toString()}
			} else if nextToken.Lexeme == "-" || nextToken.Lexeme == "!" {
				return Unary{string(lexeme), parser.parse_unary().toString()}
			} else {
				return Unary{string(lexeme), parser.parse_literal().toString()}
			}
		} else {
			return Unary{string(lexeme), nextToken.Value}
		}
	default:
		return Unary{}
	}
}

func (parser *Parser) parse_binary() Binary {
	lexeme := parser.peek().Lexeme
	parser.Current--
	left := parser.parse().toString()
	parser.advance()
	return Binary{left, string(lexeme), parser.parse().toString()}
}
