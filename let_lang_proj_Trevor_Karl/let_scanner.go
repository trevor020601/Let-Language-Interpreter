package main

import (
	"bufio"
	"unicode"
	"unicode/utf8"
)

type tokenType rune
type charClass rune

// Struct for a Token
type Token struct {
	tokenType  tokenType
	tokenValue string
}

// Struct for Character Types
type CharType struct {
	char      string
	charClass charClass
}

// Constant values for tokenTypes and charClasses
const (
	EOF               = -1
	letter  charClass = 0
	digit   charClass = 1
	unknown charClass = 99

	LParen     tokenType = 0
	RParen     tokenType = 1
	comma      tokenType = 2
	minus      tokenType = 3
	equals     tokenType = 4
	iszero     tokenType = 5
	ifKey      tokenType = 6
	then       tokenType = 7
	elseKey    tokenType = 8
	let        tokenType = 9
	in         tokenType = 10
	identifier tokenType = 11
	integer    tokenType = 12
)

// Lookup function to get tokenTypes from input
func lookup(s string) tokenType {
	switch s {
	case "(":
		return LParen
	case ")":
		return RParen
	case ",":
		return comma
	case "=":
		return equals
	default:
		return EOF
	}
}

// Separate lookup function for key words to avoid errors with lexing
func keyLookup(s string) tokenType {
	switch s {
	case "let":
		return let
	case "in":
		return in
	case "else":
		return elseKey
	case "then":
		return then
	case "if":
		return ifKey
	case "iszero":
		return iszero
	case "minus":
		return minus
	default:
		return identifier
	}
}

// Function to get the CharType of an input
func getCharType(char string) CharType {
	if char == "" {
		return CharType{char, EOF}
	}

	ch, _ := utf8.DecodeRuneInString(char)

	if unicode.IsDigit(ch) {
		return CharType{char, digit}
	}
	if unicode.IsLetter(ch) {
		return CharType{char, letter}
	}

	return CharType{char, unknown}
}

// Function to get the next char that isn't a blank
func getNextNonBlank(scanner *bufio.Scanner) string {
	input := scanner.Scan()
	if input == false {
		return ""
	}

	text := scanner.Text()
	decodedText, _ := utf8.DecodeRuneInString(text)

	for {
		if unicode.IsSpace(decodedText) == false {
			break
		}

		input = scanner.Scan()
		if input == false {
			return ""
		}

		text = scanner.Text()
		decodedText, _ = utf8.DecodeRuneInString(text)
	}

	return text
}

// Struct to create the Lexer
type Lexer struct {
	tokenQueue       []Token
	currentTokenType tokenType
	inputBuffer      *bufio.Scanner
}

// Function to create a Lexer "object"
func initializer(scanner *bufio.Scanner) *Lexer {
	return &Lexer{inputBuffer: scanner}
}

// Function for the lexical process and putting tokens into the tokenQueue
func (lex *Lexer) Lex() {
	lexeme := ""
	isEOF := false

	currentChar := getCharType(getNextNonBlank(lex.inputBuffer))
	if currentChar.char == "" {
		lex.currentTokenType = EOF
		return
	}

	if currentChar.charClass == digit {
		lexeme = lexeme + currentChar.char
		currentChar = getCharType(getNextNonBlank(lex.inputBuffer))

		for {
			if currentChar.charClass != digit {
				break
			}

			lexeme = lexeme + currentChar.char
			currentChar = getCharType(getNextNonBlank(lex.inputBuffer))
			if currentChar.char == "" {
				lex.currentTokenType = EOF
				isEOF = true
				break
			}
		}

		lex.currentTokenType = integer
		lex.tokenQueue = append(lex.tokenQueue, Token{lex.currentTokenType, lexeme})
		lexeme = ""
		if isEOF == true {
			return
		}
	}

	if currentChar.charClass == letter {
		lexeme = lexeme + currentChar.char
		lex.currentTokenType = keyLookup(lexeme)
		if lex.currentTokenType != identifier {
			lex.tokenQueue = append(lex.tokenQueue, Token{lex.currentTokenType, lexeme})
			lexeme = ""
			return
		}

		currentChar = getCharType(getNextNonBlank(lex.inputBuffer))
		if currentChar.char == "" {
			lex.currentTokenType = EOF
			return
		}

		for {
			if currentChar.charClass != letter {
				break
			}

			lexeme = lexeme + currentChar.char
			lex.currentTokenType = keyLookup(lexeme)
			if lex.currentTokenType != identifier {
				lex.tokenQueue = append(lex.tokenQueue, Token{lex.currentTokenType, lexeme})
				lexeme = ""
				return
			}

			currentChar = getCharType(getNextNonBlank(lex.inputBuffer))
			if currentChar.char == "" {
				lex.currentTokenType = EOF
				isEOF = true
				break
			}
		}

		lex.currentTokenType = keyLookup(lexeme)
		lex.tokenQueue = append(lex.tokenQueue, Token{lex.currentTokenType, lexeme})
		if isEOF == true {
			return
		}
	}

	if currentChar.charClass == unknown {
		lex.currentTokenType = lookup(currentChar.char)
		lex.tokenQueue = append(lex.tokenQueue, Token{lex.currentTokenType, currentChar.char})
	}
}
