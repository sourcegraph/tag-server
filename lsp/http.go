package lsp

import "net/http"

type HTTPLangService struct {
	svc LangSvc
}

func NewHTTPLangService(svc LangSvc) *HTTPLangService {
	return &HTTPLangService{svc: svc}
}

func (s *HTTPLangService) Initialize(req *http.Request, params *InitializeParams, result *InitializeResult) error {
	return s.svc.Initialize(params, result)
}

func (s *HTTPLangService) Completion(req *http.Request, params *TextDocumentPositionParams, result *CompletionList) error {
	return s.svc.Completion(params, result)
}

func (s *HTTPLangService) CompletionItemResolve(req *http.Request, params *CompletionList, result *CompletionList) error {
	return s.svc.CompletionItemResolve(params, result)
}

func (s *HTTPLangService) Hover(req *http.Request, params *TextDocumentPositionParams, result *Hover) error {
	return s.svc.Hover(params, result)
}

func (s *HTTPLangService) SignatureHelpRequest(req *http.Request, params *TextDocumentPositionParams, result *SignatureHelp) error {
	return s.svc.SignatureHelpRequest(params, result)
}

func (s *HTTPLangService) GoToDefinition(req *http.Request, params *TextDocumentPositionParams, result *[]Location) error {
	return s.svc.GoToDefinition(params, result)
}

func (s *HTTPLangService) References(req *http.Request, params *ReferenceParams, result *[]Location) error {
	return s.svc.References(params, result)
}

func (s *HTTPLangService) DocumentHighlights(req *http.Request, params *ReferenceParams, result *DocumentHighlight) error {
	return s.svc.DocumentHighlights(params, result)
}

func (s *HTTPLangService) DocumentSymbols(req *http.Request, params *DocumentSymbolParams, result *[]SymbolInformation) error {
	return s.svc.DocumentSymbols(params, result)
}

func (s *HTTPLangService) WorkplaceSymbols(req *http.Request, params *WorkplaceSymbolParams, result *[]SymbolInformation) error {
	return s.svc.WorkplaceSymbols(params, result)
}

func (s *HTTPLangService) CodeAction(req *http.Request, params *CodeActionParams, result *[]Command) error {
	return s.svc.CodeAction(params, result)
}

func (s *HTTPLangService) CodeLensRequest(req *http.Request, params *CodeLensParams, result *[]Command) error {
	return s.svc.CodeLensRequest(params, result)
}

func (s *HTTPLangService) CodeLensResolve(req *http.Request, params *CodeLens, result *CodeLens) error {
	return s.svc.CodeLensResolve(params, result)
}

func (s *HTTPLangService) DocumentFormatting(req *http.Request, params *DocumentFormattingParams, result *[]TextEdit) error {
	return s.svc.DocumentFormatting(params, result)
}

func (s *HTTPLangService) DocumentOnTypeFormatting(req *http.Request, params *DocumentFormattingParams, result *[]TextEdit) error {
	return s.svc.DocumentOnTypeFormatting(params, result)
}

func (s *HTTPLangService) Rename(req *http.Request, params *RenameParams, result *WorkspaceEdit) error {
	return s.svc.Rename(params, result)
}
