package lsp

import "net/http"

type LangSvc interface {
	Initialize(req *http.Request, params *InitializeParams, result *InitializeResult) error
	Completion(req *http.Request, params TextDocumentPositionParams, result CompletionList) error
	CompletionItemResolve(req *http.Request, params CompletionList, result CompletionList) error
	HoverRequest(req *http.Request, params TextDocumentPositionParams, result Hover) error
	SignatureHelpRequest(req *http.Request, params TextDocumentPositionParams, result SignatureHelp) error
	GoToDefinition(req *http.Request, params TextDocumentPositionParams, result []Location) error
	FindReferences(req *http.Request, params ReferenceParams, result []Location) error
	DocumentHighlights(req *http.Request, params ReferenceParams, result DocumentHighlight) error
	DocumentSymbols(req *http.Request, params DocumentSymbolParams, result []SymbolInformation) error
	WorkplaceSymbols(req *http.Request, params WorkplaceSymbolParams, result []SymbolInformation) error
	CodeAction(req *http.Request, params CodeActionParams, result []Command) error
	CodeLensRequest(req *http.Request, params CodeLensParams, result []Command) error
	CodeLensResolve(req *http.Request, params CodeLens, result CodeLens) error
	DocumentFormatting(req *http.Request, params DocumentFormattingParams, result []TextEdit) error
	DocumentOnTypeFormatting(req *http.Request, params DocumentFormattingParams, result []TextEdit) error
	Rename(req *http.Request, params RenameParams, result WorkspaceEdit) error
}

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

func Initialize(req *http.Request, params *InitializeParams, result *InitializeResult) error {
	return nil
}

var TextDocumentSyncKind = map[string]int32{
	/**
	* Defines how the host (editor) should sync document changes to the language server.
	 */
	"None": 0,

	/**
	 * Documents are synced by always sending the full content of the document.
	 */
	"Full": 1,

	/**
	 * Documents are synced by sending the full content on open. After that only incremental
	 * updates to the document are sent.
	 */
	"Incremental": 2,
}

type ServerCapabilities struct {
	/**
	 * Defines how text documents are synced.
	 */
	TextDocumentSync int32
	/**
	 * The server provides hover support.
	 */
	HoverProvider bool
	/**
	 * The server provides completion support.
	 */
	CompletionProvider CompletionOptions
	/**
	 * The server provides signature help support.
	 */
	SignatureHelpProvider SignatureHelpOptions
	/**
	 * The server provides goto definition support.
	 */
	DefinitionProvider bool
	/**
	 * The server provides find references support.
	 */
	ReferencesProvider bool
	/**
	 * The server provides document highlight support.
	 */
	DocumentHighlightProvider bool
	/**
	 * The server provides document symbol support.
	 */
	DocumentSymbolProvider bool
	/**
	 * The server provides workspace symbol support.
	 */
	WorkspaceSymbolProvider bool
	/**
	 * The server provides code actions.
	 */
	CodeActionProvider bool
	/**
	 * The server provides code lens.
	 */
	CodeLensProvider CodeLensOptions
	/**
	 * The server provides document formatting.
	 */
	DocumentFormattingProvider bool
	/**
	 * The server provides document range formatting.
	 */
	DocumentRangeFormattingProvider bool
	/**
	 * The server provides document formatting on typing.
	 */
	DocumentOnTypeFormattingProvider DocumentOnTypeFormattingOptions
	/**
	 * The server provides rename support.
	 */
	RenameProvider bool
}

/**
 * Completion options.
 */
type CompletionOptions struct {
	/**
	 * The server provides support to resolve additional information for a completion item.
	 */
	ResolveProvider bool

	/**
	 * The characters that trigger completion automatically.
	 */
	TriggerCharacters []string
}

/**
 * Format document on type options
 */
type DocumentOnTypeFormattingOptions struct {
	/**
	 * A character on which formatting should be triggered, like `}`.
	 */
	FirstTriggerCharacter string

	/**
	 * More trigger characters.
	 */
	MoreTriggerCharacter []string
}

type CodeLensOptions struct {
	/**
	 * Code lens has a resolve provider as well.
	 */
	ResolveProvider bool
}

type SignatureHelpOptions struct {
	TriggerCharacters []string
}

var CompletionItemKind = map[string]int32{
	"Text":        1,
	"Method":      2,
	"Function":    3,
	"Constructor": 4,
	"Field":       5,
	"Variable":    6,
	"Class":       7,
	"Interface":   8,
	"Module":      9,
	"Property":    10,
	"Unit":        11,
	"Value":       12,
	"Enum":        13,
	"Keyword":     14,
	"Snippet":     15,
	"Color":       16,
	"File":        17,
	"Reference":   18,
}

type CompletionItem struct {
	/**
	 * The label of this completion item. By default
	 * also the text that is inserted when selecting
	 * this completion.
	 */
	Label string
	/**
	 * The kind of this completion item. Based of the kind
	 * an icon is chosen by the editor.
	 */
	Kind int32
	/**
	 * A human-readable string with additional information
	 * about this item, like type or symbol information.
	 */
	Detail string
	/**
	 * A human-readable string that represents a doc-comment.
	 */
	Documentation string
	/**
	 * A string that shoud be used when comparing this item
	 * with other items. When `falsy` the label is used.
	 */
	SortText string
	/**
	 * A string that should be used when filtering a set of
	 * completion items. When `falsy` the label is used.
	 */
	FilterText string
	/**
	 * A string that should be inserted a document when selecting
	 * this completion. When `falsy` the label is used.
	 */
	InsertText string
	/**
	 * An edit which is applied to a document when selecting
	 * this completion. When an edit is provided the value of
	 * insertText is ignored.
	 */
	TextEdit TextEdit
	/**
	 * An data entry field that is preserved on a completion item between
	 * a completion and a completion resolve request.
	 */
	Data interface{}
}

/**
 * Represents a collection of [completion items](#CompletionItem) to be presented
 * in the editor.
 */
type CompletionList struct {
	/**
	 * This list it not complete. Further typing should result in recomputing
	 * this list.
	 */
	IsIncomplete bool

	/**
	 * The completion items.
	 */
	Items []CompletionItem
}

func Completion(req *http.Request, params TextDocumentPositionParams, result CompletionList) error {
	return nil
}

func CompletionItemResolve(req *http.Request, params CompletionList, result CompletionList) error {
	return nil
}

type Hover struct {
	/**
	 * The hover's content
	 */
	Contents []MarkedString

	/**
	 * An optional range
	 */
	Range Range
}

type MarkedString struct {
	Language string

	Value string
}

func HoverRequest(req *http.Request, params TextDocumentPositionParams, result Hover) error {
	return nil
}

/**
 * Signature help represents the signature of something
 * callable. There can be multiple signature but only one
 * active and only one active parameter.
 */
type SignatureHelp struct {
	/**
	 * One or more signatures.
	 */
	Signatures []SignatureInformation

	/**
	 * The active signature.
	 */
	ActiveSignature int32

	/**
	 * The active parameter of the active signature.
	 */
	ActiveParameter int32
}

/**
 * Represents the signature of something callable. A signature
 * can have a label, like a function-name, a doc-comment, and
 * a set of parameters.
 */
type SignatureInformation struct {
	/**
	 * The label of this signature. Will be shown in
	 * the UI.
	 */
	Label string

	/**
	 * The human-readable doc-comment of this signature. Will be shown
	 * in the UI but can be omitted.
	 */
	Documentation string

	/**
	 * The parameters of this signature.
	 */
	Paramaters []ParameterInformation
}

/**
 * Represents a parameter of a callable-signature. A parameter can
 * have a label and a doc-comment.
 */
type ParameterInformation struct {
	/**
	 * The label of this signature. Will be shown in
	 * the UI.
	 */
	Label string

	/**
	 * The human-readable doc-comment of this signature. Will be shown
	 * in the UI but can be omitted.
	 */
	Documentation string
}

func SignatureHelpRequest(req *http.Request, params TextDocumentPositionParams, result SignatureHelp) error {
	return nil
}

func GoToDefinition(req *http.Request, params TextDocumentPositionParams, result []Location) error {
	return nil
}

type ReferenceContext struct {
	/**
	 * Include the declaration of the current symbol.
	 */
	IncludeDeclaration bool
}

type ReferenceParams struct {
	TextDocumentPositionParams

	Context ReferenceContext
}

func FindReferences(req *http.Request, params ReferenceParams, result []Location) error {
	return nil
}

/**
 * A document highlight kind.
 */
var DocumentHighlightKind = map[string]int32{
	/**
	 * A textual occurrance.
	 */
	"Text": 1,

	/**
	 * Read-access of a symbol, like reading a variable.
	 */
	"Read": 2,

	/**
	 * Write-access of a symbol, like writing to a variable.
	 */
	"Write": 3,
}

/**
 * A document highlight is a range inside a text document which deserves
 * special attention. Usually a document highlight is visualized by changing
 * the background color of its range.
 */
type DocumentHighlight struct {
	/**
	 * The range this highlight applies to.
	 */
	Range Range

	/**
	 * The highlight kind, default is DocumentHighlightKind.Text.
	 */
	Kind int32
}

func DocumentHighlights(req *http.Request, params ReferenceParams, result DocumentHighlight) error {
	return nil
}

type DocumentSymbolParams struct {
	/**
	 * The text document.
	 */
	TextDocument TextDocumentIdentifier
}

var SymbolKind = map[string]int32{
	"File":        1,
	"Module":      2,
	"Namespace":   3,
	"Package":     4,
	"Class":       5,
	"Method":      6,
	"Property":    7,
	"Field":       8,
	"Constructor": 9,
	"Enum":        10,
	"Interface":   11,
	"Function":    12,
	"Variable":    13,
	"Constant":    14,
	"String":      15,
	"Number":      16,
	"Boolean":     17,
	"Array":       18,
}

/**
 * Represents information about programming constructs like variables, classes,
 * interfaces etc.
 */
type SymbolInformation struct {
	/**
	 * The name of this symbol.
	 */
	Name string

	/**
	 * The kind of this symbol.
	 */
	Kind int32

	/**
	 * The location of this symbol.
	 */
	Location Location

	/**
	 * The name of the symbol containing this symbol.
	 */
	ContainerName string
}

func DocumentSymbols(req *http.Request, params DocumentSymbolParams, result []SymbolInformation) error {
	return nil
}

/**
 * The parameters of a Workspace Symbol Request.
 */
type WorkplaceSymbolParams struct {
	/**
	 * A non-empty query string
	 */
	Query string
}

func WorkplaceSymbols(req *http.Request, params WorkplaceSymbolParams, result []SymbolInformation) error {
	return nil
}

/**
 * Contains additional diagnostic information about the context in which
 * a code action is run.
 */
type CodeActionContext struct {
	/**
	 * An array of diagnostics.
	 */
	Diagnostics []Diagnostic
}

/**
 * Params for the CodeActionRequest
 */
type CodeActionParams struct {
	/**
	 * The document in which the command was invoked.
	 */
	TextDocument TextDocumentIdentifier

	/**
	 * The range for which the command was invoked.
	 */
	Range Range

	/**
	 * Context carrying additional information.
	 */
	Context CodeActionContext
}

func CodeAction(req *http.Request, params CodeActionParams, result []Command) error {
	return nil
}

type CodeLensParams struct {
	/**
	 * The document to request code lens for.
	 */
	TextDocument TextDocumentIdentifier
}

/**
 * A code lens represents a command that should be shown along with
 * source text, like the number of references, a way to run tests, etc.
 *
 * A code lens is _unresolved_ when no command is associated to it. For performance
 * reasons the creation of a code lens and resolving should be done in two stages.
 */
type CodeLens struct {
	/**
	 * The range in which this code lens is valid. Should only span a single line.
	 */
	Range Range

	/**
	 * The command this code lens represents.
	 */
	Command Command

	/**
	 * A data entry field that is preserved on a code lens item between
	 * a code lens and a code lens resolve request.
	 */
	Data interface{}
}

func CodeLensRequest(req *http.Request, params CodeLensParams, result []Command) error {
	return nil
}

func CodeLensResolve(req *http.Request, params CodeLens, result CodeLens) error {
	return nil
}

type DocumentFormattingParams struct {
	/**
	 * The document to format.
	 */
	TextDocument TextDocumentIdentifier

	/**
	 * The format options.
	 */
	Options FormattingOptions
}

/**
 * Value-object describing what options formatting should use.
 */
type FormattingOptions struct {
	/**
	 * Size of a tab in spaces.
	 */
	TabSize int32

	/**
	 * Prefer spaces over tabs.
	 */
	InsertSpaces bool

	/**
	* Signature for further properites
	 */
	Key string
}

func DocumentFormatting(req *http.Request, params DocumentFormattingParams, result []TextEdit) error {
	return nil
}

func DocumentOnTypeFormatting(req *http.Request, params DocumentFormattingParams, result []TextEdit) error {
	return nil
}

type RenameParams struct {
	/**
	 * The document to format.
	 */
	TextDocument TextDocumentIdentifier

	/**
	 * The position at which this request was sent.
	 */
	Position Position

	/**
	 * The new name of the symbol. If the given name is not valid the
	 * request must return a [ResponseError](#ResponseError) with an
	 * appropriate message set.
	 */
	NewName string
}

func Rename(req *http.Request, params RenameParams, result WorkspaceEdit) error {
	return nil
}
