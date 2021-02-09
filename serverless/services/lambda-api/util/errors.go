package util

import "fmt"

type ErrorDataNotFoundForSymbol struct {
	DataType string
	Symbol   string
}

func (e *ErrorDataNotFoundForSymbol) Error() string {
	return fmt.Sprintf("ERROR: %s Data not found for %s", e.DataType, e.Symbol)
}

func NewErrorDataNotFoundForSymbol(dataType string, symbol string) *ErrorDataNotFoundForSymbol {
	return &ErrorDataNotFoundForSymbol{DataType: dataType, Symbol: symbol}
}

type ErrorMethodNotImplemented struct {
	Method string
}

func (e *ErrorMethodNotImplemented) Error() string {
	return fmt.Sprintf("ERROR: '%s' method not implemented for this route", e.Method)
}

func NewErrorMethodNotImplemented(method string) *ErrorMethodNotImplemented {
	return &ErrorMethodNotImplemented{Method: method}
}
