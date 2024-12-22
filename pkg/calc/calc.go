package calc

import (
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	Unknown TokenType = iota
	Number
	Plus
	Minus
	Multiply
	Divide
	LeftParen
	RightParen
)

// Token представляет собой токен с типом и значением
type Token struct {
	Type  TokenType
	Value string
}

// Tokenize разбивает строку-выражение на токены
func Tokenize(expression string) ([]Token, error) {
	var tokens []Token
	var current string

	for _, r := range expression {
		switch {
		case unicode.IsDigit(r):
			current += string(r)
		case r == '+':
			if current != "" {
				tokens = append(tokens, Token{Number, current})
				current = ""
			}
			tokens = append(tokens, Token{Plus, "+"})
		case r == '-':
			if current != "" {
				tokens = append(tokens, Token{Number, current})
				current = ""
			}
			tokens = append(tokens, Token{Minus, "-"})
		case r == '*':
			if current != "" {
				tokens = append(tokens, Token{Number, current})
				current = ""
			}
			tokens = append(tokens, Token{Multiply, "*"})
		case r == '/':
			if current != "" {
				tokens = append(tokens, Token{Number, current})
				current = ""
			}
			tokens = append(tokens, Token{Divide, "/"})
		case r == '(':
			if current != "" {
				return nil, ErrInvalidCharacter
			}
			tokens = append(tokens, Token{LeftParen, "("})
		case r == ')':
			if current != "" {
				tokens = append(tokens, Token{Number, current})
				current = ""
			}
			tokens = append(tokens, Token{RightParen, ")"})
		default:
			return nil, ErrInvalidCharacter
		}
	}

	if current != "" {
		tokens = append(tokens, Token{Number, current})
	}

	return tokens, nil
}

var precedence = map[TokenType]int{
	Plus:     1,
	Minus:    1,
	Multiply: 2,
	Divide:   2,
}

func ShuntingYard(tokens []Token) ([]Token, error) {
	var output []Token
	var operators []Token

	for _, token := range tokens {
		switch token.Type {
		case Number:
			output = append(output, token)
		case Plus, Minus, Multiply, Divide:
			for len(operators) > 0 && precedence[operators[len(operators)-1].Type] >= precedence[token.Type] {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		case LeftParen:
			operators = append(operators, token)
		case RightParen:
			for len(operators) > 0 && operators[len(operators)-1].Type != LeftParen {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, ErrMismatchedParentheses
			}
			operators = operators[:len(operators)-1]
		default:
			return nil, ErrUnexpectedToken
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1].Type == LeftParen {
			return nil, ErrMismatchedParentheses
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// EvaluateRPN вычисляет выражение в обратной польской нотации
func EvaluateRPN(tokens []Token) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token.Type {
		case Number:
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, value)
		case Plus:
			if len(stack) < 2 {
				return 0, ErrNotEnoughValues
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a+b)
		case Minus:
			if len(stack) < 2 {
				return 0, ErrNotEnoughValues
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a-b)
		case Multiply:
			if len(stack) < 2 {
				return 0, ErrNotEnoughValues
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a*b)
		case Divide:
			if len(stack) < 2 {
				return 0, ErrNotEnoughValues
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			if b == 0 {
				return 0, ErrDivisionByZero
			}
			stack = stack[:len(stack)-2]
			stack = append(stack, a/b)
		default:
			return 0, ErrUnexpectedToken
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, ErrEmptyInput
	}
	expression = strings.ReplaceAll(expression, " ", "")
	tokens, err := Tokenize(expression)
	if err != nil {
		return 0, err
	}

	rpn, err := ShuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	result, err := EvaluateRPN(rpn)
	if err != nil {
		return 0, err
	}
	return result, nil
}
