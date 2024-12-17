package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
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
	contentsLength := len(fileContents)
	exitCode := 0
	line := 1

	for idx := 0; idx < contentsLength; idx++ {

		switch fileContents[idx] {
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
			if idx+1 < contentsLength && fileContents[idx+1] == '=' {
				fmt.Println("EQUAL_EQUAL == null")
				idx++
			} else {
				fmt.Println("EQUAL = null")
			}
		case '!':
			if idx+1 < contentsLength && fileContents[idx+1] == '=' {
				fmt.Println("BANG_EQUAL != null")
				idx++
			} else {
				fmt.Println("BANG ! null")
			}
		case '<':
			if idx+1 < contentsLength && fileContents[idx+1] == '=' {
				fmt.Println("LESS_EQUAL <= null")
				idx++
			} else {
				fmt.Println("LESS < null")
			}
		case '>':
			if idx+1 < contentsLength && fileContents[idx+1] == '=' {
				fmt.Println("GREATER_EQUAL >= null")
				idx++
			} else {
				fmt.Println("GREATER > null")
			}
		case '/':
			if idx+1 < contentsLength && fileContents[idx+1] == '/' {
				for idx < contentsLength && fileContents[idx] != '\n' {
					idx++
				}
				line++
			} else {
				fmt.Println("SLASH / null")
			}
		case '\n':
			line++
			continue
		case '\t':
			continue
		case '\r':
			continue
		case ' ':
			continue
		case '"':
			var bytes []byte
			for idx < contentsLength {
				idx++
				if idx == contentsLength {
					fmt.Fprintf(os.Stderr, "[line %v] Error: Unterminated string.", line)
					exitCode = 65
				} else if fileContents[idx] == '"' {
					str := string(bytes)
					fmt.Printf("STRING \"%v\" %v\n", str, str)
					break
				} else {
					bytes = append(bytes, fileContents[idx])
				}
			}
		default:
			if isNumerical(fileContents[idx]) {
				var numericalBytes []byte

				for idx+1 < contentsLength && isNumerical(fileContents[idx+1]) {
					numericalBytes = append(numericalBytes, fileContents[idx])
					idx++
				}

				numericalBytes = append(numericalBytes, fileContents[idx])
				stringLiteral := string(numericalBytes)
				floatingLiteral, _ := strconv.ParseFloat(stringLiteral, 64)
				if strings.Contains(stringLiteral, ".") {
					fmt.Printf("NUMBER %v %.1v\n", floatingLiteral, strconv.FormatFloat(floatingLiteral, 'f', -1, 64))
				} else {
					fmt.Printf("NUMBER %v %.1f\n", floatingLiteral, floatingLiteral)
				}
			} else {
				fmt.Fprintf(os.Stderr, "[line %v] Error: Unexpected character: %v\n", line, string(fileContents[idx]))
				exitCode = 65
			}
		}
	}

	fmt.Println("EOF  null")
	os.Exit(exitCode)
}

func isNumerical(b byte) bool {
	return b >= 48 && b <= 57 || b == '.'
}
