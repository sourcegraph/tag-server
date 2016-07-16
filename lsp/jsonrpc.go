package lsp

import "fmt"

type ResponseError struct {
	/**
	 * A number indicating the error type that occurred.
	 */
	Code int

	/**
	 * A string providing a short description of the error.
	 */
	Message string

	/**
	 * A Primitive or Structured value that contains additional
	 * information about the error. Can be omitted.
	 */
	Data interface{}
}

type ErrorCode int

const (
	ParseError       ErrorCode = -32700
	InvalidRequest             = -32600
	MethodNotFound             = -32601
	InvalidParams              = -32602
	InternalError              = -32603
	serverErrorStart           = -3209
	serverErrorEnd             = -32000
)

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%+v", e)
}
