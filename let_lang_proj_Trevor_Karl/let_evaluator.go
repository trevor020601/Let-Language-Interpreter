package main

import (
	"os"
	"strconv"
)

// Looks up the name of a binding and returns the binding value
func bindLookup(b []Binding, name string) string {
	str := ""
	for i := range b {
		if b[i].varname == name {
			str = b[i].value
		}
	}
	return str
}

// Struct for the Evaluator object
type Evaluator struct {
	root *astNode
}

// Initializes the evaluator with a passed root node
func evalInitializer(r *astNode) Evaluator {
	return Evaluator{r}
}

// Helper function to call the evaluate process
func (e *Evaluator) eval() string {
	return e.evaluate(e.root, []Binding{})
}

// Function for the evaluation process
func (e *Evaluator) evaluate(localRoot *astNode, b []Binding) string {
	localRoot.env = b

	if localRoot.ttype == let {
		varname := localRoot.children[0].tvalue
		exp1Val := e.evaluate(localRoot.children[1], b)
		b = append(b, Binding{varname, exp1Val})
		return e.evaluate(localRoot.children[2], b)
	}
	if localRoot.ttype == minus {
		exp1Val, err := strconv.Atoi(e.evaluate(localRoot.children[0], b))
		if err != nil {
			os.Exit(1)
		}
		exp2Val, err := strconv.Atoi(e.evaluate(localRoot.children[1], b))
		if err != nil {
			os.Exit(1)
		}
		return strconv.Itoa(exp1Val - exp2Val)
	}
	if localRoot.ttype == iszero {
		expVal, err := strconv.Atoi(e.evaluate(localRoot.children[0], b))
		if err != nil {
			os.Exit(1)
		}
		return strconv.FormatBool(expVal == 0)
	}
	if localRoot.ttype == ifKey {
		expBoolVal := e.evaluate(localRoot.children[0], b)
		if expBoolVal == "true" {
			return e.evaluate(localRoot.children[1], b)
		}
		return e.evaluate(localRoot.children[2], b)
	}
	if localRoot.ttype == identifier {
		return bindLookup(b, localRoot.tvalue)
	}
	if localRoot.ttype == integer {
		return localRoot.tvalue
	}
	return ""
}
