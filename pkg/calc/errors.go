package calc

import "errors"

var (
	ErrInvalidExpression     = errors.New("invalid expression")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrEmptyInput            = errors.New("empty input")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
	ErrUnexpectedToken       = errors.New("unexpected token")
	ErrNotEnoughValues       = errors.New("not enough values in expression")
	ErrInvalidCharacter      = errors.New("invalid character")
)
