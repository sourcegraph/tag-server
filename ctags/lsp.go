package ctags

import (
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"

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
	log.Printf("GoToDefinition(%+v)", params)

	docURL, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(docURL.Path)
	if err != nil {
		return err
	}
	file := string(b)
	lines := strings.Split(file, "\n")
	line := lines[params.Position.Line]
	start := strings.LastIndexAny(line[:params.Position.Character], " \r\n\t()\"'.,*-") + 1
	end := strings.IndexAny(line[params.Position.Character:], " \r\n\t()\"'.,*-")
	if end == -1 {
		end = len(line)
	} else {
		end += params.Position.Character
	}
	token := line[start:end] // This is the token to search for

	log.Printf("search around for token %q", token)

	var matchedTags []Tag
	{
		searchDir := filepath.Dir(docURL.Path)
		dirfiles, err := ioutil.ReadDir(searchDir)
		if err != nil {
			return err
		}
		for _, file := range dirfiles {
			if file.IsDir() {
				continue
			}
			parser, err := Parse([]string{filepath.Join(searchDir, file.Name())})
			if err != nil {
				return err
			}
			for _, tag := range parser.Tags() {
				if tag.Name == token {
					matchedTags = append(matchedTags, tag)
				}
			}
		}
	}

	log.Printf("matched %d tags", len(matchedTags))

	symbols := tagsToSymbolInformation(matchedTags)
	locs := make([]lsp.Location, len(symbols))
	for i, symbol := range symbols {
		locs[i] = symbol.Location
	}

	*result = locs
	return nil
}
func (s *LangSvc) References(params *lsp.ReferenceParams, result *[]lsp.Location) error {
	log.Printf("References(%+v)", params)

	return nil
}
func (s *LangSvc) DocumentHighlights(params *lsp.ReferenceParams, result *lsp.DocumentHighlight) error {
	return nil
}
func (s *LangSvc) DocumentSymbols(params *lsp.DocumentSymbolParams, result *[]lsp.SymbolInformation) error {
	docURL, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return err
	}

	parser, err := Parse([]string{docURL.Path})
	if err != nil {
		return err
	}
	*result = tagsToSymbolInformation(parser.Tags())
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

func tagsToSymbolInformation(tags []Tag) []lsp.SymbolInformation {
	res := make([]lsp.SymbolInformation, 0, len(tags))
	for _, tag := range tags {
		nameIdx := strings.Index(tag.Def, tag.Name)
		res = append(res, lsp.SymbolInformation{
			Name: tag.Name,
			Kind: lsp.SKMethod, // TODO
			Location: lsp.Location{
				URI: "file://" + tag.File,
				Range: lsp.Range{
					Start: lsp.Position{Line: tag.Line - 1, Character: nameIdx},
					End:   lsp.Position{Line: tag.Line - 1, Character: nameIdx + len(tag.Name)},
				},
			},
		})
	}
	return res
}
