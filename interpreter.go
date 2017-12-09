package calcgo

import (
	"errors"
	"strconv"
)

// Errors, that can be returned by the interpreter
var (
	ErrorMissingLeftChild  = errors.New("Error: Missing left child of node")
	ErrorMissingRightChild = errors.New("Error: Missing right child of node")
	ErrorInvalidNodeType   = errors.New("Error: Invalid node type")
	ErrorInvalidInteger    = errors.New("Error: Invalid Integer")
	ErrorInvalidDecimal    = errors.New("Error: Invalid Decimal")
	ErrorInvalidVariable   = errors.New("Error: Invalid Variable")
	ErrorParserError       = errors.New("Error: Parser error")
	ErrorDivisionByZero    = errors.New("Error: Division by zero")
)

// Interpret interprets a given string.
// Can return an error if parsing failed
//
// Examples:
//  caclgo.Interpret("(1 + 2) * 3") // Result: 9
//  caclgo.Interpret("1 + 2 * 3")   // Result: 7
func Interpret(str string) (float64, []error) {
	if len(str) == 0 {
		return 0, nil
	}

	ast, errors := Parse(str)
	if errors != nil {
		return 0, errors
	}

	result, err := InterpretAST(ast)
	if err != nil {
		return 0, []error{err}
	}

	return result, nil
}

// InterpretAST interprets a given ast.
// Can return an error if the ast is invalid.
func InterpretAST(ast AST) (float64, error) {
	return calculateNode(ast.Node)
}

func calculateNode(node *Node) (float64, error) {
	switch node.Type {
	case NInteger:
		integer, err := strconv.Atoi(node.Value)
		if err != nil {
			return 0, ErrorInvalidInteger
		}
		return float64(integer), nil
	case NDecimal:
		decimal, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return 0, ErrorInvalidDecimal
		}
		return decimal, nil
	}

	if !IsOperator(node.Type) {
		return 0, ErrorInvalidNodeType
	}

	if node.LeftChild == nil {
		return 0, ErrorMissingLeftChild
	}
	if node.RightChild == nil {
		return 0, ErrorMissingRightChild
	}

	left, err := calculateNode(node.LeftChild)
	if err != nil {
		return 0, err
	}
	right, err := calculateNode(node.RightChild)
	if err != nil {
		return 0, err
	}

	var result float64
	switch node.Type {
	case NAddition:
		result = left + right
	case NSubtraction:
		result = left - right
	case NMultiplication:
		result = left * right
	case NDivision:
		if right == 0 {
			return 0, ErrorDivisionByZero
		}
		result = left / right
	}

	return result, nil
}
