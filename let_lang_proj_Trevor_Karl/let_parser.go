package main

import (
	"fmt"
	"os"
)

// Binding represents a var/value pair
type Binding struct {
	varname string
	value   string //rune
}

// Prints a Binding as a string
func printBind(bind Binding) string {
	return "(" + bind.varname + ", " + bind.value + ")"
}

// Struct for an abstract tree node
type astNode struct {
	parent   *astNode  // pointer to parent node if exists. If null, then root.
	ttype    tokenType // token type, used to distinguish btw vars and integers.
	termsym  bool      // is it a terminal symbol (leaf node)
	tvalue   string    // represent the contents as a string, even if it's an int
	children []*astNode
	env      []Binding
}

// Converts tokenTypes to a string for printing purposes
func printToken(tt tokenType) string {
	switch tt {
	case identifier:
		return "VarExp"
	case integer:
		return "ConstExp"
	case minus:
		return "DiffExp"
	case iszero:
		return "IsZeroExp"
	case ifKey:
		return "IfExp"
	case let:
		return "LetExp"
	default:
		return ""
	}
}

// Prints the abstract syntax tree WITH environments
func (node *astNode) printASTWithEnv(numIndents int) {
	str := ""
	needsParen := false

	for i := 0; i < numIndents; i++ {
		str = str + "    "
	}

	fmt.Print(str)
	fmt.Print(printToken(node.ttype))

	for r, bind := range node.env {
		if r == 0 {
			fmt.Print(" -> Env [")
		}
		if r < len(node.env) && r > 0 {
			fmt.Print(", ")
		}

		fmt.Print(printBind(bind))

		if r == len(node.env)-1 {
			fmt.Print("]")
		}
	}

	if node.ttype == identifier || node.ttype == integer || node.ttype == minus ||
		node.ttype == iszero || node.ttype == ifKey || node.ttype == let {
		needsParen = true
		fmt.Print(" (\n")
	}

	if len(node.children) > 0 && node.termsym == false {
		for c := 0; c < len(node.children); c++ {
			node.children[c].printASTWithEnv(numIndents + 1)
			//Fix for minus and let
			if (node.ttype == minus && c == 0) ||
				(node.ttype == let && c == 1) {
				fmt.Print(",")
			}
			fmt.Print("\n")
		}
	}

	identifierStr := ""
	for s := 0; s < len(printToken(node.ttype))-2; s++ {
		identifierStr = identifierStr + " "
	}
	if node.ttype == identifier {
		fmt.Print(str + identifierStr)
		//Put double quotes around identifiers
		fmt.Printf("\"%s\"\n", node.tvalue)
	}

	if node.ttype == integer {
		fmt.Print(str + identifierStr)
		fmt.Printf(node.tvalue + "\n")
	}

	fmt.Print(str)
	if needsParen == true {
		fmt.Print(")")
	}
}

// Calls printASTWithEnv with the root node
func printTreeWE(nd *astNode) {
	nd.printASTWithEnv(0)
}

// Slices the first Token off the queue
func popTQ(tQ *[]Token) Token {
	queue := *tQ
	tok := queue[0]
	queue = queue[1:]
	*tQ = queue
	return tok
}

// Creates an object for the Parser
type Parser struct {
	root *astNode
	tQ   []Token
}

// Pops the token queue to move to the next token
func (par *Parser) advanceToken() {
	popTQ(&par.tQ)
}

// Peeks at the next value in the token queue and returns it
func (par *Parser) peekTQ() (tokenType, string) {
	if len(par.tQ) <= 0 {
		return tokenType(unknown), ""
	}

	return par.tQ[0].tokenType, par.tQ[0].tokenValue
}

// Checks next expected token and advances the queue if it is correct
func (par *Parser) expected(tt tokenType, advance bool, err string) {
	nextT, _ := par.peekTQ()

	if tt != nextT {
		fmt.Println(err)
		os.Exit(1)
	}

	if advance == true {
		par.advanceToken()
	}
}

// Initialer for tree nodes
func (par *Parser) initTreeNd(nd *astNode) {
	nd.termsym = true
	nd.children = make([]*astNode, 0, 3) // 5
	nd.ttype, nd.tvalue = par.peekTQ()
	par.advanceToken()
}

// Feeds the token queue into the parser to begin parsing process
func parseInput(tQ []Token) *astNode {
	par := Parser{&astNode{}, tQ}
	par.parse()
	return par.root
}

// Function for the parsing process and creating nodes
func (par *Parser) parse() astNode {
	parent := astNode{}

	if par.root.children == nil {
		par.root = &parent
	}

	par.initTreeNd(&parent)

	if parent.ttype == minus {
		parent.termsym = false
		par.expected(LParen, true, "Expected Left Parenthesis, Unexpected Token Received Instead")
		left := par.parse()
		par.expected(comma, true, "Expected Comma, Unexpected Token Received Instead")
		right := par.parse()
		par.expected(RParen, true, "Expected Right Parenthesis, Unexpected Token Received Instead")
		left.parent = &parent
		parent.children = append(parent.children, &left)
		right.parent = &parent
		parent.children = append(parent.children, &right)
	}
	if parent.ttype == iszero {
		parent.termsym = false
		par.expected(LParen, true, "Expected Left Parenthesis, Unexpected Token Received Instead")
		expressionChild := par.parse()
		par.expected(RParen, true, "Expected Right Parenthesis, Unexpected Token Received Instead")
		expressionChild.parent = &parent
		parent.children = append(parent.children, &expressionChild)
	}
	if parent.ttype == ifKey {
		parent.termsym = false
		par.expected(iszero, false, "Expected Iszero, Unexpected Token Received Instead")
		predicate := par.parse()
		par.expected(then, true, "Expected Then, Unexpected Token Received Instead")
		falseCase := par.parse()
		par.expected(elseKey, true, "Expected Else, Unexpected Token Received Instead")
		trueCase := par.parse()
		predicate.parent = &parent
		parent.children = append(parent.children, &predicate)
		falseCase.parent = &parent
		parent.children = append(parent.children, &falseCase)
		trueCase.parent = &parent
		parent.children = append(parent.children, &trueCase)
	}
	if parent.ttype == let {
		parent.termsym = false
		par.expected(identifier, false, "Expected Identifier, Unexpected Token Received Instead")
		id := par.parse()
		par.expected(equals, true, "Expected Equals, Unexpected Token Received Instead")
		firstChild := par.parse()
		par.expected(in, true, "Expected In, Unexpected Token Received Instead")
		secondChild := par.parse()
		id.parent = &parent
		parent.children = append(parent.children, &id)
		firstChild.parent = &parent
		parent.children = append(parent.children, &firstChild)
		secondChild.parent = &parent
		parent.children = append(parent.children, &secondChild)
	}
	return parent
}
