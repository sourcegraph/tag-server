package server

import (
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/tag-server/ctags"

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

	file, err := fetchFile(params.TextDocument.URI)
	if err != nil {
		return err
	}
	token, _ := extractTokenFromPosition(file, params.Position.Line, params.Position.Character) // token to search for

	log.Printf("search around for token %q", token)

	var matchedTags []ctags.Tag
	{
		docURL, err := url.Parse(params.TextDocument.URI)
		if err != nil {
			return err
		}
		searchDir := filepath.Dir(docURL.Path)
		dirfiles, err := ioutil.ReadDir(searchDir)
		if err != nil {
			return err
		}
		for _, file := range dirfiles {
			if file.IsDir() {
				continue
			}
			parser, err := ctags.Parse2([]string{filepath.Join(searchDir, file.Name())})
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

	file, err := fetchFile(params.TextDocument.URI)
	if err != nil {
		return err
	}
	token, _ := extractTokenFromPosition(file, params.Position.Line, params.Position.Character)

	var matchedRefs []lsp.Location
	{
		docURL, err := url.Parse(params.TextDocument.URI)
		if err != nil {
			return err
		}
		searchDir := filepath.Dir(docURL.Path)
		dirfiles, err := ioutil.ReadDir(searchDir)
		if err != nil {
			return err
		}
		for _, file := range dirfiles {
			if file.IsDir() {
				continue
			}
			b, err := ioutil.ReadFile(filepath.Join(searchDir, file.Name()))
			if err != nil {
				return err
			}
			lines := strings.Split(string(b), "\n")

			for l, line := range lines {
				if c := strings.Index(line, token); c != -1 {
					matchedRefs = append(matchedRefs, lsp.Location{
						URI: "file://" + filepath.Join(searchDir, file.Name()),
						Range: lsp.Range{
							Start: lsp.Position{Line: l, Character: c},
							End:   lsp.Position{Line: l, Character: c + len(token)},
						},
					})
				}
			}
		}
	}

	*result = matchedRefs
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

	parser, err := ctags.Parse2([]string{docURL.Path})
	if err != nil {
		return err
	}
	*result = tagsToSymbolInformation(parser.Tags())
	return nil
}
func (s *LangSvc) WorkspaceSymbols(params *lsp.WorkspaceSymbolParams, result *[]lsp.SymbolInformation) error {
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

var nameToSymbolKind = map[string]lsp.SymbolKind{
	"file":        lsp.SKFile,
	"module":      lsp.SKModule,
	"namespace":   lsp.SKNamespace,
	"package":     lsp.SKPackage,
	"class":       lsp.SKClass,
	"method":      lsp.SKMethod,
	"property":    lsp.SKProperty,
	"field":       lsp.SKField,
	"constructor": lsp.SKConstructor,
	"enum":        lsp.SKEnum,
	"interface":   lsp.SKInterface,
	"function":    lsp.SKFunction,
	"variable":    lsp.SKVariable,
	"constant":    lsp.SKConstant,
	"string":      lsp.SKString,
	"number":      lsp.SKNumber,
	"boolean":     lsp.SKBoolean,
	"array":       lsp.SKArray,
}

func tagsToSymbolInformation(tags []ctags.Tag) []lsp.SymbolInformation {
	res := make([]lsp.SymbolInformation, 0, len(tags))
	for _, tag := range tags {
		nameIdx := strings.Index(tag.DefLinePrefix, tag.Name)
		if nameIdx < 0 {
			log.Printf("! dropping tag because could not find name (%s) in def line prefix (%q)", tag.Name, tag.DefLinePrefix)
			continue
		}
		kind := nameToSymbolKind[tag.Kind]
		if kind == 0 {
			kind = lsp.SKVariable
		}
		res = append(res, lsp.SymbolInformation{
			Name: tag.Name,
			Kind: kind,
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

func etagsToSymbolInformation(tags []ctags.ETag) []lsp.SymbolInformation {
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

func extractTokenFromPosition(file string, l int, c int) (token string, loc lsp.Range) {
	lines := strings.Split(file, "\n")
	line := lines[l]
	start := strings.LastIndexAny(line[:c], " \r\n\t()\"'.,*-<>:") + 1
	end := strings.IndexAny(line[c:], " \r\n\t()\"'.,*-<>:")
	if end == -1 {
		end = len(line)
	} else {
		end += c
	}
	return line[start:end], lsp.Range{
		Start: lsp.Position{Line: l, Character: start},
		End:   lsp.Position{Line: l, Character: end},
	}
}

// fetches file contents from URI
func fetchFile(uri string) (string, error) {
	docURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(docURL.Path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
