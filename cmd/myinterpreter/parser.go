package main

import (
	"fmt"
)

type Literal struct {
	Value string
}

type Group struct {
	Expr Statement
}

type Unary struct {
	Operator string
	Expr     Statement
}

type Binary struct {
	Left     Statement
	Operator string
	Right    Statement
}

type Statement interface {
	toString() string
}

func (literal Literal) toString() string {
	return literal.Value
}

func (group Group) toString() string {
	return fmt.Sprintf("(group %s)", group.Expr.toString())
}

func (unary Unary) toString() string {
	return fmt.Sprintf("(%s %s)", unary.Operator, unary.Expr.toString())
}

func (binary Binary) toString() string {
	return fmt.Sprintf("(%s %s %s)", binary.Left.toString(), binary.Operator, binary.Right.toString())
}

type Parser struct {
	Tokens  []Token
	Current int
}

func (parser *Parser) peek() Token {
	if parser.Current < len(parser.Tokens) {
		return parser.Tokens[parser.Current]
	}
	return Token{Kind: "EOF"}
}

func (parser *Parser) advance() Token {
	if !parser.atTheEnd() {
		parser.Current++
	}
	return parser.previous()
}

func (parser *Parser) atTheEnd() bool {
	return parser.Current >= len(parser.Tokens)
}

func (parser *Parser) previous() Token {
	if parser.Current > 0 {
		return parser.Tokens[parser.Current-1]
	}
	return Token{Kind: "EOF"}
}

func (parser *Parser) parse() (Statement, error) {
	return parser.parseExpression()
}

func (parser *Parser) parseExpression() (Statement, error) {
	return parser.parseTerm()
}

func (parser *Parser) parseTerm() (Statement, error) {
	expr, err := parser.parseFactor()
	if err != nil {
		return nil, err
	}

	for parser.peek().Kind == "*" || parser.peek().Kind == "/" {
		token := parser.advance()
		right, err := parser.parseFactor()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: string(token.Lexeme), Right: right}
	}

	return expr, nil
}

func (parser *Parser) parseFactor() (Statement, error) {
	switch parser.peek().Kind {
	case "(":
		return parser.parseGroup()
	case "-", "!":
		return parser.parseUnary()
	case "NUMBER", "STRING", "TRUE", "FALSE", "NIL":
		return parser.parseLiteral(), nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", parser.peek())
	}
}

func (parser *Parser) parseLiteral() Literal {
	token := parser.advance()
	return Literal{Value: token.Value}
}

func (parser *Parser) parseGroup() (Group, error) {
	parser.advance() // Consume '('

	expr, err := parser.parseExpression()
	if err != nil {
		return Group{}, err
	}

	if parser.peek().Kind != ")" {
		return Group{}, fmt.Errorf("expected ')' but got %v", parser.peek())
	}
	parser.advance() // Consume ')'

	return Group{Expr: expr}, nil
}

func (parser *Parser) parseUnary() (Unary, error) {
	token := parser.advance()

	expr, err := parser.parseFactor()
	if err != nil {
		return Unary{}, err
	}

	return Unary{Operator: string(token.Lexeme), Expr: expr}, nil
}

func (parser *Parser) parseBinary(left Statement) (Binary, error) {
	operator := parser.advance().Lexeme
	right, err := parser.parseFactor()
	if err != nil {
		return Binary{}, err
	}

	return Binary{Left: left, Operator: string(operator), Right: right}, nil
}
