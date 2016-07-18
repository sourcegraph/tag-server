package ctags

import (
	"log"

	"github.com/sourcegraph/tag-server/lsp"
)

type langSvc struct{}

func (s *langSvc) Initialize(params *lsp.InitializeParams, result *lsp.InitializeResult) error {
	log.Printf("langSvc.Initialize(%+v)", params)
	// vfsURL, err := url.Parse(params.RootPath)
	// if err != nil {
	// 	return fmt.Errorf("Could not parse VFS URL: %s", err)
	// }

	// vfsURL := "http://localhost:7979"

	return nil
}
func (s *langSvc) Completion(params *lsp.TextDocumentPositionParams, result *lsp.CompletionList) error {
	return nil
}
func (s *langSvc) CompletionItemResolve(params *lsp.CompletionList, result *lsp.CompletionList) error {
	return nil
}
func (s *langSvc) HoverRequest(params *lsp.TextDocumentPositionParams, result *lsp.Hover) error {
	return nil
}
func (s *langSvc) SignatureHelpRequest(params *lsp.TextDocumentPositionParams, result *lsp.SignatureHelp) error {
	return nil
}
func (s *langSvc) GoToDefinition(params *lsp.TextDocumentPositionParams, result *[]lsp.Location) error {
	// TODO
	return nil
}
func (s *langSvc) FindReferences(params *lsp.ReferenceParams, result *[]lsp.Location) error {
	return nil
}
func (s *langSvc) DocumentHighlights(params *lsp.ReferenceParams, result *lsp.DocumentHighlight) error {
	return nil
}
func (s *langSvc) DocumentSymbols(params *lsp.DocumentSymbolParams, result *[]lsp.SymbolInformation) error {
	// TODO
	return nil
}
func (s *langSvc) WorkplaceSymbols(params *lsp.WorkplaceSymbolParams, result *[]lsp.SymbolInformation) error {
	return nil
}
func (s *langSvc) CodeAction(params *lsp.CodeActionParams, result *[]lsp.Command) error {
	return nil
}
func (s *langSvc) CodeLensRequest(params *lsp.CodeLensParams, result *[]lsp.Command) error {
	return nil
}
func (s *langSvc) CodeLensResolve(params *lsp.CodeLens, result *lsp.CodeLens) error {
	return nil
}
func (s *langSvc) DocumentFormatting(params *lsp.DocumentFormattingParams, result *[]lsp.TextEdit) error {
	return nil
}
func (s *langSvc) DocumentOnTypeFormatting(params *lsp.DocumentFormattingParams, result *[]lsp.TextEdit) error {
	return nil
}
func (s *langSvc) Rename(params *lsp.RenameParams, result *lsp.WorkspaceEdit) error {
	return nil
}

var _ lsp.LangSvc = (*langSvc)(nil)

func NewLangService() lsp.LangSvc {
	return &langSvc{}
}
