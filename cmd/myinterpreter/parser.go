package main

import (
	"fmt"
	"os"
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

func (p *Parser) parse() Statement {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Statement {
	return p.parseComparison()
}

func (p *Parser) parseComparison() Statement {
	expr := p.parseTerm()

	for p.match("LESS", "LESS_EQUAL", "GREATER", "GREATER_EQUAL", "EQUAL_EQUAL", "BANG_EQUAL") {
		operator := p.previous()
		right := p.parseTerm()
		expr = Binary{Left: expr, Operator: string(operator.Lexeme), Right: right}
	}

	return expr
}

func (p *Parser) parseTerm() Statement {
	expr := p.parseFactor()

	for p.match("PLUS", "MINUS") {
		operator := p.previous()
		right := p.parseFactor()
		expr = Binary{Left: expr, Operator: string(operator.Lexeme), Right: right}
	}

	return expr
}

func (p *Parser) parseFactor() Statement {
	expr := p.parseUnary()

	for p.match("STAR", "SLASH") {
		operator := p.previous()
		right := p.parseUnary()
		expr = Binary{Left: expr, Operator: string(operator.Lexeme), Right: right}
	}

	return expr
}

func (p *Parser) parseUnary() Statement {
	if p.match("BANG", "MINUS") {
		operator := p.previous()
		right := p.parseUnary()
		return Unary{Operator: string(operator.Lexeme), Expr: right}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() Statement {
	if p.match("TRUE", "FALSE", "NIL") {
		return Literal{Value: string(p.previous().Lexeme)}
	}

	if p.match("NUMBER", "STRING") {
		return Literal{Value: p.previous().Value}
	}

	if p.match("LEFT_PAREN") {
		expr := p.parseExpression()
		if !p.match("RIGHT_PAREN") {
			fmt.Fprintf(os.Stderr, "[line %v] Error at '%v': Expect expression.", line, p.peek().Lexeme)
			errorCode = 65
		}
		return Group{Expr: expr}
	}

	fmt.Fprintf(os.Stderr, "[line %v] Error at '%v': Expect expression.", line, p.peek().Lexeme)
	errorCode = 65

	return Literal{}
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
