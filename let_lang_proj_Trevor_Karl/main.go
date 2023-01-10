package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	commandLine := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter .let example to evaluate. Preferably example.let or if_example.let:")
	examp := ""

	for commandLine.Scan() {
		if _, e := os.Stat(commandLine.Text()); os.IsNotExist(e) == false {
			examp = commandLine.Text()
			break
		} else {
			fmt.Println("ERROR - File does not exist or is not in current directory")
		}
	}

	exampPath := examp

	exampInput, e := os.ReadFile(exampPath)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	exampContents := string(exampInput)
	contents := bufio.NewScanner(strings.NewReader(exampContents))
	contents.Split(bufio.ScanRunes)

	fmt.Printf("\nContents of %s:\n", examp)
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(exampContents)

	lexer := initializer(contents)
	for {
		if lexer.currentTokenType == EOF {
			lexer.tokenQueue = append(lexer.tokenQueue, Token{EOF, "EOF"})
			break
		}

		lexer.Lex()
	}

	fmt.Println()
	fmt.Println("Token Queue: ")
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println()
	for _, index := range lexer.tokenQueue {
		fmt.Println("Token:", index.tokenType, "   Lexeme:", index.tokenValue)
	}

	root := parseInput(lexer.tokenQueue)
	fmt.Println()
	fmt.Println("Abstract Syntax Tree Without Environments")
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println()
	printTreeWE(root)

	fmt.Println()

	evaluator := evalInitializer(root)
	fmt.Println("\nEvaluation")
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("\nExpression evaluated to: " + evaluator.eval())

	fmt.Println()

	fmt.Println("Abstract Syntax Tree With Environments")
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println()
	printTreeWE(evaluator.root)
}
