package ctags

import (
	"log"

	"sourcegraph.com/sourcegraph/sourcegraph/pkg/lsp"
)

type LangSvc struct {
	RootPath string
}

var Server = &LangSvc{}

func (s *LangSvc) Initialize(params *lsp.InitializeParams, result *lsp.InitializeResult) error {
	log.Printf("LangSvc.Initialize(%+v)", params)
	log.Printf("root path: %q", params.RootPath)
	s.RootPath = params.RootPath
	result.Capabilities = lsp.ServerCapabilities{
		HoverProvider:          true,
		DocumentSymbolProvider: true,
		DefinitionProvider:     true,
		ReferencesProvider:     true,
	}

	return nil
}
func (s *LangSvc) Completion(params *lsp.TextDocumentPositionParams, result *lsp.CompletionList) error {
	return nil
}
func (s *LangSvc) CompletionItemResolve(params *lsp.CompletionList, result *lsp.CompletionList) error {
	return nil
}
func (s *LangSvc) Hover(params *lsp.TextDocumentPositionParams, result *lsp.Hover) error {
	result.Contents = []lsp.MarkedString{{Language: "markdown", Value: "Hello CTags!"}}
	return nil
}
func (s *LangSvc) SignatureHelpRequest(params *lsp.TextDocumentPositionParams, result *lsp.SignatureHelp) error {
	return nil
}
func (s *LangSvc) GoToDefinition(params *lsp.TextDocumentPositionParams, result *[]lsp.Location) error {
	// TODO
	return nil
}
func (s *LangSvc) References(params *lsp.ReferenceParams, result *[]lsp.Location) error {
	return nil
}
func (s *LangSvc) DocumentHighlights(params *lsp.ReferenceParams, result *lsp.DocumentHighlight) error {
	return nil
}
func (s *LangSvc) DocumentSymbols(params *lsp.DocumentSymbolParams, result *[]lsp.SymbolInformation) error {
	// TODO
	return nil
}
func (s *LangSvc) WorkplaceSymbols(params *lsp.WorkplaceSymbolParams, result *[]lsp.SymbolInformation) error {
	return nil
}
func (s *LangSvc) CodeAction(params *lsp.CodeActionParams, result *[]lsp.Command) error {
	return nil
}
func (s *LangSvc) CodeLensRequest(params *lsp.CodeLensParams, result *[]lsp.Command) error {
	return nil
}
func (s *LangSvc) CodeLensResolve(params *lsp.CodeLens, result *lsp.CodeLens) error {
	return nil
}
func (s *LangSvc) DocumentFormatting(params *lsp.DocumentFormattingParams, result *[]lsp.TextEdit) error {
	return nil
}
func (s *LangSvc) DocumentOnTypeFormatting(params *lsp.DocumentFormattingParams, result *[]lsp.TextEdit) error {
	return nil
}
func (s *LangSvc) Rename(params *lsp.RenameParams, result *lsp.WorkspaceEdit) error {
	return nil
}
