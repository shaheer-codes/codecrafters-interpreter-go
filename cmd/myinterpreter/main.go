package main

import (
	"fmt"
	"os"
)

// const (
// 	LEFT_PAREN  rune = '('
// 	RIGHT_PAREN rune = ')'
// 	LEFT_BRACE  rune = '{'
// 	RIGHT_BRACE rune = '}'
// 	COMMA       rune = ','
// 	DOT         rune = '.'
// 	MINUS       rune = '-'
// 	PLUS        rune = '+'
// 	SEMICOLON   rune = ';'
// 	STAR        rune = '*'
// )

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	rawFileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fileContents := string(rawFileContents)
	errFlag := false
	for _, curr := range fileContents {
		switch curr {
		case '(':
			fmt.Println("LEFT_PAREN ( null")
		case ')':
			fmt.Println("RIGHT_PAREN ) null")
		case '{':
			fmt.Println("LEFT_BRACE { null")
		case '}':
			fmt.Println("RIGHT_BRACE } null")
		case ',':
			fmt.Println("COMMA , null")
		case '.':
			fmt.Println("DOT . null")
		case '-':
			fmt.Println("MINUS - null")
		case '+':
			fmt.Println("PLUS + null")
		case ';':
			fmt.Println("SEMICOLON ; null")
		case '*':
			fmt.Println("STAR * null")
		case '$':
			fmt.Fprintln(os.Stderr, "[line 1] Error: Unexpected character: $")
			errFlag = true
			os.Exit(65)
		case '#':
			fmt.Fprintln(os.Stderr, "[line 1] Error: Unexpected character: #")
			errFlag = true
			os.Exit(65)
		case '@':
			fmt.Fprintln(os.Stderr, "[line 1] Error: Unexpected character: @")
			errFlag = true
			os.Exit(65)
		case '~':
			fmt.Fprintln(os.Stderr, "[line 1] Error: Unexpected character: ~")
			errFlag = true
		}
	}

	fmt.Println("EOF  null")
	fmt.Fprintln(os.Stderr, "EOF  null")
	if errFlag {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}
