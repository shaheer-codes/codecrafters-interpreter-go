package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TokenType string

type Token struct {
	Kind   string
	Lexeme TokenType
	Value  string
}

const (
	LEFT_PAREN    TokenType = "("
	RIGHT_PAREN   TokenType = ")"
	LEFT_BRACE    TokenType = "{"
	RIGHT_BRACE   TokenType = "}"
	COMMA         TokenType = ","
	DOT           TokenType = "."
	MINUS         TokenType = "-"
	PLUS          TokenType = "+"
	SEMICOLON     TokenType = ";"
	STAR          TokenType = "*"
	EQUAL         TokenType = "="
	EQUAL_EQUAL   TokenType = "=="
	BANG          TokenType = "!"
	BANG_EQUAL    TokenType = "!="
	LESS          TokenType = "<"
	LESS_EQUAL    TokenType = "<="
	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="
	SLASH         TokenType = "/"
	EOF           TokenType = ""
	AND           TokenType = "AND"
	CLASS         TokenType = "CLASS"
	ELSE          TokenType = "ELSE"
	FALSE         TokenType = "FALSE"
	FOR           TokenType = "FOR"
	FUN           TokenType = "FUN"
	IF            TokenType = "IF"
	NIL           TokenType = "NIL"
	OR            TokenType = "OR"
	PRINT         TokenType = "PRINT"
	RETURN        TokenType = "RETURN"
	SUPER         TokenType = "SUPER"
	THIS          TokenType = "THIS"
	TRUE          TokenType = "TRUE"
	VAR           TokenType = "VAR"
	WHILE         TokenType = "WHILE"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

var line = 1
var errorCode = 0

func main() {
	fmt.Fprintln(os.Stderr, "Logs from the program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./interpreter.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	switch command {
	case "tokenize":
		tokens := tokenize(filename)
		for i := 0; i < len(tokens); i++ {
			fmt.Printf("%v %v %v\n", tokens[i].Kind, tokens[i].Lexeme, tokens[i].Value)
		}
	case "parse":
		tokens := tokenize(filename)

		parser := Parser{tokens, 0}

		for parser.peek().Kind != "EOF" {
			statement, err := parser.parse()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				errorCode = 65
			} else {
				fmt.Printf("%v", statement.toString())
			}

			parser.advance()
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	os.Exit(errorCode)
}

func NewToken(name string, lexeme TokenType, value string) Token {
	return Token{Kind: name, Lexeme: lexeme, Value: value}
}

type Lexer struct {
	Input string
	Pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{Input: input, Pos: 0}
}

func tokenize(filename string) []Token {
	rawFileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fileContents := string(rawFileContents)

	lexer := NewLexer(fileContents)

	var tokens []Token

	for {
		token := lexer.nextToken()
		tokens = append(tokens, token)
		if token.Kind == "EOF" {
			break
		}
	}

	return tokens
}

func (lexer *Lexer) peek() byte {
	if lexer.Pos >= len(lexer.Input) {
		return 0
	}

	return lexer.Input[lexer.Pos]
}

func (lexer *Lexer) readByte() byte {
	if lexer.Pos >= len(lexer.Input) {
		return 0
	}

	ch := lexer.Input[lexer.Pos]
	lexer.Pos++

	return ch
}

func (lexer *Lexer) peekNext() byte {
	if lexer.Pos+1 >= len(lexer.Input) {
		return 0
	}

	return lexer.Input[lexer.Pos+1]
}

func (lexer *Lexer) atTheEnd() bool {
	return lexer.Pos == len(lexer.Input)
}

func (lexer *Lexer) skipWhitespaces() {
	for lexer.peek() == '\n' || lexer.peek() == '\t' || lexer.peek() == ' ' || lexer.peek() == '\r' {
		if lexer.peek() == '\n' {
			line++
		}
		lexer.readByte()
	}
}

func (lexer *Lexer) nextToken() Token {
	var token Token

	lexer.skipWhitespaces()

	switch lexer.peek() {
	case '(':
		token = NewToken("LEFT_PAREN", LEFT_PAREN, "null")
	case ')':
		token = NewToken("RIGHT_PAREN", RIGHT_PAREN, "null")
	case '{':
		token = NewToken("LEFT_BRACE", LEFT_BRACE, "null")
	case '}':
		token = NewToken("RIGHT_BRACE", RIGHT_BRACE, "null")
	case ',':
		token = NewToken("COMMA", COMMA, "null")
	case '.':
		token = NewToken("DOT", DOT, "null")
	case '-':
		token = NewToken("MINUS", MINUS, "null")
	case '+':
		token = NewToken("PLUS", PLUS, "null")
	case ';':
		token = NewToken("SEMICOLON", SEMICOLON, "null")
	case '*':
		token = NewToken("STAR", STAR, "null")
	case '=':
		if lexer.peekNext() == '=' {
			token = NewToken("EQUAL_EQUAL", EQUAL_EQUAL, "null")
			lexer.readByte()
		} else {
			token = NewToken("EQUAL", EQUAL, "null")
		}
	case '!':
		if lexer.peekNext() == '=' {
			token = NewToken("BANG_EQUAL", BANG_EQUAL, "null")
			lexer.readByte()
		} else {
			token = NewToken("BANG", BANG, "null")
		}
	case '<':
		if lexer.peekNext() == '=' {
			token = NewToken("LESS_EQUAL", LESS_EQUAL, "null")
			lexer.readByte()
		} else {
			token = NewToken("LESS", LESS, "null")
		}
	case '>':
		if lexer.peekNext() == '=' {
			token = NewToken("GREATER_EQUAL", GREATER_EQUAL, "null")
			lexer.readByte()
		} else {
			token = NewToken("GREATER", GREATER, "null")
		}
	case '/':
		if lexer.peekNext() == '/' {
			for lexer.peek() != '\n' && !lexer.atTheEnd() {
				lexer.readByte()
			}

			return lexer.nextToken()
		} else {
			token = NewToken("SLASH", SLASH, "null")
		}
	case '"':
		start := lexer.Pos + 1
		lexer.readByte()
		for lexer.peek() != '"' && !lexer.atTheEnd() {
			lexer.readByte()
		}
		if lexer.atTheEnd() {
			fmt.Fprintf(os.Stderr, "[line %v] Error: Unterminated string.\n", line)
			errorCode = 65
			return lexer.nextToken()
		} else {
			stringValue := lexer.Input[start:lexer.Pos]
			token = NewToken("STRING", TokenType(fmt.Sprintf("\"%v\"", stringValue)), stringValue)
		}
	case 0:
		token = NewToken("EOF", EOF, "null")
	default:
		if unicode.IsDigit(rune(lexer.peek())) {
			start := lexer.Pos
			for unicode.IsDigit(rune(lexer.peek())) {
				lexer.readByte()
			}

			if lexer.peek() == '.' && unicode.IsDigit(rune(lexer.peekNext())) {
				lexer.readByte()
				for unicode.IsDigit(rune(lexer.peek())) {
					lexer.readByte()
				}
			}

			number := lexer.Input[start:lexer.Pos]

			floatNumber, _ := strconv.ParseFloat(number, 64)

			if !strings.Contains(number, ".") || !strings.Contains(fmt.Sprintf("%v", floatNumber), ".") {
				return NewToken("NUMBER", TokenType(number), fmt.Sprintf("%.1f", floatNumber))
			} else {
				return NewToken("NUMBER", TokenType(number), strconv.FormatFloat(floatNumber, 'f', -1, 64))
			}
		} else if unicode.IsLetter(rune(lexer.peek())) || lexer.peek() == '_' {
			start := lexer.Pos
			for unicode.IsLetter(rune(lexer.peek())) || unicode.IsDigit(rune(lexer.peek())) || lexer.peek() == '_' {
				lexer.readByte()
			}

			identifier := lexer.Input[start:lexer.Pos]

			if keywords[identifier] != "" {
				return NewToken(string(keywords[identifier]), TokenType(identifier), "null")
			} else {
				return NewToken("IDENTIFIER", TokenType(identifier), "null")
			}
		} else {
			fmt.Fprintf(os.Stderr, "[line %v] Error: Unexpected character: %v\n", line, string(lexer.peek()))
			errorCode = 65
			lexer.readByte()
			return lexer.nextToken()
		}
	}

	lexer.readByte()
	return token
}
