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

	exitCode := 0

	var storedEqual, storedBang bool

	for _, curr := range fileContents {

		if curr != '=' && storedEqual {
			fmt.Println("EQUAL = null")
			storedEqual = false
		}

		if curr != '=' && storedBang {
			fmt.Println("BANG ! null")
			storedBang = false
		}

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
		case '=':
			if storedEqual {
				fmt.Println("EQUAL_EQUAL == null")
				storedEqual = false
			} else if storedBang {
				fmt.Println("BANG_EQUAL != null")
				storedBang = false
			} else {
				storedEqual = true
			}
		case '!':
			storedBang = true
		default:
			fmt.Fprintf(os.Stderr, "[line 1] Error: Unexpected character: %v\n", string(curr))
			exitCode = 65
		}
	}

	if storedEqual {
		fmt.Println("EQUAL = null")
		storedEqual = false
	}

	if storedBang {
		fmt.Println("BANG ! null")
		storedBang = false
	}

	fmt.Println("EOF  null")
	os.Exit(exitCode)
}
