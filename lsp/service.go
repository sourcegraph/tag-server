package lsp

import "net/http"

type LangSvc struct{}

type InitializeParams struct {
	/**
	 * The process Id of the parent process that started
	 * the server.
	 */
	ProcessId int

	/**
	 * The rootPath of the workspace. Is null
	 * if no folder is open.
	 */
	RootPath string

	/**
	 * The capabilities provided by the client (editor)
	 */
	Capabilities ClientCapabilities
}

type ClientCapabilities struct{}

type InitializeResult struct {
	/**
	 * The capabilities the language server provides.
	 */
	Capabilities ServerCapabilities
}

type InitializeError struct {
	/**
	 * Indicates whether the client should retry to send the
	 * initilize request after showing the message provided
	 * in the ResponseError.
	 */
	Retry bool
}

func (s *LangSvc) Initialize(req *http.Request, params *InitializeParams, result *InitializeResult) error {
	return nil
}

type ServerCapabilities struct {
	// TODO
}
