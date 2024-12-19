package main

import (
	"fmt"
)

type Statement interface {
	toString() string
}

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
	return fmt.Sprintf("(%s %s %s)", binary.Operator, binary.Left.toString(), binary.Right.toString())
}

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) peek() Token {
	if p.Current < len(p.Tokens) {
		return p.Tokens[p.Current]
	}
	return Token{Kind: "EOF"}
}

func (p *Parser) advance() Token {
	if !p.atTheEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) previous() Token {
	if p.Current > 0 {
		return p.Tokens[p.Current-1]
	}
	return Token{Kind: "EOF"}
}

func (p *Parser) atTheEnd() bool {
	return p.Current >= len(p.Tokens)
}

func (p *Parser) parse() (Statement, error) {
	return p.parseExpression()
}

func (p *Parser) parseExpression() (Statement, error) {
	return p.parseComparison()
}

func (p *Parser) parseComparison() (Statement, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return Binary{}, err
	}

	for p.match("LESS", "LESS_EQUAL", "GREATER", "GREATER_EQUAL", "EQUAL_EQUAL", "BANG_EQUAL") {
		operator := p.previous()
		right, err := p.parseTerm()
		if err != nil {
			return Binary{}, err
		}
		expr = Binary{Left: expr, Operator: string(operator.Lexeme), Right: right}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Statement, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return Binary{}, err
	}

	for p.match("PLUS", "MINUS") {
		operator := p.previous()
		right, err := p.parseFactor()
		if err != nil {
			return Binary{}, err
		}
		expr = Binary{Left: expr, Operator: string(operator.Lexeme), Right: right}
	}

	return expr, nil
}

func (p *Parser) parseFactor() (Statement, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return Unary{}, err
	}

	for p.match("STAR", "SLASH") {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return Binary{}, err
		}
		expr = Binary{Left: expr, Operator: string(operator.Lexeme), Right: right}
	}

	return expr, nil
}

func (p *Parser) parseUnary() (Statement, error) {
	if p.match("BANG", "MINUS") {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return Unary{}, err
		}
		return Unary{Operator: string(operator.Lexeme), Expr: right}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Statement, error) {
	if p.match("TRUE", "FALSE", "NIL") {
		return Literal{Value: string(p.previous().Lexeme)}, nil
	}

	if p.match("NUMBER", "STRING") {
		return Literal{Value: p.previous().Value}, nil
	}

	if p.match("LEFT_PAREN") {
		expr, err := p.parseExpression()
		if err != nil {
			return Group{}, err
		}
		if !p.match("RIGHT_PAREN") {
			return Group{}, fmt.Errorf("[line %v] Error: Expected )", line)
		}
		return Group{Expr: expr}, fmt.Errorf("[line %v] Error at %v: Expected expression", line, p.peek())
	}

	panic(fmt.Sprintf("Unexpected token: %v", p.peek()))
}

func (p *Parser) match(types ...string) bool {
	if p.atTheEnd() {
		return false
	}

	for _, t := range types {
		if p.peek().Kind == t {
			p.advance()
			return true
		}
	}

	return false
}
