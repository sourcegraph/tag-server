package lsp

type LangSvc interface {
	Initialize(params *InitializeParams, result *InitializeResult) error
	Completion(params *TextDocumentPositionParams, result *CompletionList) error
	CompletionItemResolve(params *CompletionList, result *CompletionList) error
	HoverRequest(params *TextDocumentPositionParams, result *Hover) error
	SignatureHelpRequest(params *TextDocumentPositionParams, result *SignatureHelp) error
	GoToDefinition(params *TextDocumentPositionParams, result *[]Location) error
	FindReferences(params *ReferenceParams, result *[]Location) error
	DocumentHighlights(params *ReferenceParams, result *DocumentHighlight) error
	DocumentSymbols(params *DocumentSymbolParams, result *[]SymbolInformation) error
	WorkplaceSymbols(params *WorkplaceSymbolParams, result *[]SymbolInformation) error
	CodeAction(params *CodeActionParams, result *[]Command) error
	CodeLensRequest(params *CodeLensParams, result *[]Command) error
	CodeLensResolve(params *CodeLens, result *CodeLens) error
	DocumentFormatting(params *DocumentFormattingParams, result *[]TextEdit) error
	DocumentOnTypeFormatting(params *DocumentFormattingParams, result *[]TextEdit) error
	Rename(params *RenameParams, result *WorkspaceEdit) error
}

type None struct{}

type InitializeParams struct {
	/**
	 * The process Id of the parent process that started
	 * the server.
	 */
	ProcessID int `json:"processId"`

	/**
	 * The rootPath of the workspace. Is null
	 * if no folder is open.
	 */
	RootPath string `json:"rootPath"`

	/**
	 * The capabilities provided by the client (editor)
	 */
	Capabilities ClientCapabilities `json:"capabilities"`
}

type ClientCapabilities struct{}

type InitializeResult struct {
	/**
	 * The capabilities the language server provides.
	 */
	Capabilities ServerCapabilities `json:"capabilities"`
}

type InitializeError struct {
	/**
	 * Indicates whether the client should retry to send the
	 * initilize request after showing the message provided
	 * in the ResponseError.
	 */
	Retry bool `json:"retry"`
}

type TextDocumentSyncKind int

const (
	/**
	 * Defines how the host (editor) should sync document changes to the language server.
	 */
	TDSKNone TextDocumentSyncKind = 0

	/**
	 * Documents are synced by always sending the full content of the document.
	 */
	TDSKFull = 1

	/**
	 * Documents are synced by sending the full content on open. After that only incremental
	 * updates to the document are sent.
	 */
	TDSKIncremental = 2
)

type ServerCapabilities struct {
	/**
	 * Defines how text documents are synced.
	 */
	TextDocumentSync int32 `json:"textDocumentSync,omitempty"`
	/**
	 * The server provides hover support.
	 */
	HoverProvider bool `json:"hoverProvider,omitempty"`
	/**
	 * The server provides completion support.
	 */
	CompletionProvider CompletionOptions `json:"completionProvider,omitempty"`
	/**
	 * The server provides signature help support.
	 */
	SignatureHelpProvider SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`
	/**
	 * The server provides goto definition support.
	 */
	DefinitionProvider bool `json:"definitionProvider,omitempty"`
	/**
	 * The server provides find references support.
	 */
	ReferencesProvider bool `json:"referencesProvider,omitempty"`
	/**
	 * The server provides document highlight support.
	 */
	DocumentHighlightProvider bool `json:"documentHighlightProvider,omitempty"`
	/**
	 * The server provides document symbol support.
	 */
	DocumentSymbolProvider bool `json:"documentSymbolProvider,omitempty"`
	/**
	 * The server provides workspace symbol support.
	 */
	WorkspaceSymbolProvider bool `json:"workspaceSymbolProvider,omitempty"`
	/**
	 * The server provides code actions.
	 */
	CodeActionProvider bool `json:"codeActionProvider,omitempty"`
	/**
	 * The server provides code lens.
	 */
	CodeLensProvider CodeLensOptions `json:"codeLensProvider,omitempty"`
	/**
	 * The server provides document formatting.
	 */
	DocumentFormattingProvider bool `json:"documentFormattingProvider,omitempty"`
	/**
	 * The server provides document range formatting.
	 */
	DocumentRangeFormattingProvider bool `json:"documentRangeFormattingProvider,omitempty"`
	/**
	 * The server provides document formatting on typing.
	 */
	DocumentOnTypeFormattingProvider DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	/**
	 * The server provides rename support.
	 */
	RenameProvider bool `json:"renameProvider,omitempty"`
}

/**
 * Completion options.
 */
type CompletionOptions struct {
	/**
	 * The server provides support to resolve additional information for a completion item.
	 */
	ResolveProvider bool `json:"resolveProvider,omitempty"`

	/**
	 * The characters that trigger completion automatically.
	 */
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

/**
 * Format document on type options
 */
type DocumentOnTypeFormattingOptions struct {
	/**
	 * A character on which formatting should be triggered, like `}`.
	 */
	FirstTriggerCharacter string `json:"firstTriggerCharacter"`

	/**
	 * More trigger characters.
	 */
	MoreTriggerCharacter []string `json:"moreTriggerCharacter,omitempty"`
}

type CodeLensOptions struct {
	/**
	 * Code lens has a resolve provider as well.
	 */
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type SignatureHelpOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

type CompletionItemKind int

const (
	CIKText        CompletionItemKind = 1
	CIKMethod                         = 2
	CIKFunction                       = 3
	CIKConstructor                    = 4
	CIKField                          = 5
	CIKVariable                       = 6
	CIKClass                          = 7
	CIKInterface                      = 8
	CIKModule                         = 9
	CIKProperty                       = 10
	CIKUnit                           = 11
	CIKValue                          = 12
	CIKEnum                           = 13
	CIKKeyword                        = 14
	CIKSnippet                        = 15
	CIKColor                          = 16
	CIKFile                           = 17
	CIKReference                      = 18
)

type CompletionItem struct {
	/**
	 * The label of this completion item. By default
	 * also the text that is inserted when selecting
	 * this completion.
	 */
	Label string `json:"label"`
	/**
	 * The kind of this completion item. Based of the kind
	 * an icon is chosen by the editor.
	 */
	Kind int32 `json:"kind,omitempty"`
	/**
	 * A human-readable string with additional information
	 * about this item, like type or symbol information.
	 */
	Detail string `json:"detail,omitempty"`
	/**
	 * A human-readable string that represents a doc-comment.
	 */
	Documentation string `json:"documentation,omitempty"`
	/**
	 * A string that shoud be used when comparing this item
	 * with other items. When `falsy` the label is used.
	 */
	SortText string `json:"sortText,omitempty"`
	/**
	 * A string that should be used when filtering a set of
	 * completion items. When `falsy` the label is used.
	 */
	FilterText string `json:"filterText,omitempty"`
	/**
	 * A string that should be inserted a document when selecting
	 * this completion. When `falsy` the label is used.
	 */
	InsertText string `json:"insertText,omitempty"`
	/**
	 * An edit which is applied to a document when selecting
	 * this completion. When an edit is provided the value of
	 * insertText is ignored.
	 */
	TextEdit TextEdit `json:"textEdit,omitempty"`
	/**
	 * An data entry field that is preserved on a completion item between
	 * a completion and a completion resolve request.
	 */
	Data interface{} `json:"data,omitempty"`
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
	IsIncomplete bool `json:"isIncomplete"`

	/**
	 * The completion items.
	 */
	Items []CompletionItem `json:"items"`
}

type Hover struct {
	/**
	 * The hover's content
	 */
	Contents []MarkedString `json:"contents,omitempty"`

	/**
	 * An optional range
	 */
	Range Range `json:"range"`
}

type MarkedString struct {
	Language string `json:"language"`

	Value string `json:"value"`
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
	Signatures []SignatureInformation `json:"signatures"`

	/**
	 * The active signature.
	 */
	ActiveSignature int32 `json:"activeSignature,omitempty"`

	/**
	 * The active parameter of the active signature.
	 */
	ActiveParameter int32 `json:"activeParameter,omitempty"`
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
	Label string `json:"label"`

	/**
	 * The human-readable doc-comment of this signature. Will be shown
	 * in the UI but can be omitted.
	 */
	Documentation string `json:"documentation,omitempty"`

	/**
	 * The parameters of this signature.
	 */
	Paramaters []ParameterInformation `json:"paramaters,omitempty"`
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
	Label string `json:"label"`

	/**
	 * The human-readable doc-comment of this signature. Will be shown
	 * in the UI but can be omitted.
	 */
	Documentation string `json:"documentation,omitempty"`
}

type ReferenceContext struct {
	/**
	 * Include the declaration of the current symbol.
	 */
	IncludeDeclaration bool `json:"IncludeDeclaration"`
}

type ReferenceParams struct {
	TextDocumentPositionParams

	Context ReferenceContext `json:"context"`
}

/**
 * A document highlight kind.
 */
type DocumentHighlightKind int

const (
	/**
	 * A textual occurrance.
	 */
	Text DocumentHighlightKind = 1

	/**
	 * Read-access of a symbol, like reading a variable.
	 */
	Read = 2

	/**
	 * Write-access of a symbol, like writing to a variable.
	 */
	Write = 3
)

/**
 * A document highlight is a range inside a text document which deserves
 * special attention. Usually a document highlight is visualized by changing
 * the background color of its range.
 */
type DocumentHighlight struct {
	/**
	 * The range this highlight applies to.
	 */
	Range Range `json:"range"`

	/**
	 * The highlight kind, default is DocumentHighlightKind.Text.
	 */
	Kind int32 `json:"kind,omitempty"`
}

type DocumentSymbolParams struct {
	/**
	 * The text document.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SymbolKind int

const (
	SKFile        SymbolKind = 1
	SKModule                 = 2
	Namespace                = 3
	Package                  = 4
	SKClass                  = 5
	SKMethod                 = 6
	SKProperty               = 7
	SKField                  = 8
	SKConstructor            = 9
	SKEnum                   = 10
	SKInterface              = 11
	SKFunction               = 12
	SKVariable               = 13
	Constant                 = 14
	String                   = 15
	Number                   = 16
	Boolean                  = 17
	Array                    = 18
)

/**
 * Represents information about programming constructs like variables, classes,
 * interfaces etc.
 */
type SymbolInformation struct {
	/**
	 * The name of this symbol.
	 */
	Name string `json:"name"`

	/**
	 * The kind of this symbol.
	 */
	Kind int32 `json:"kind"`

	/**
	 * The location of this symbol.
	 */
	Location Location `json:"location"`

	/**
	 * The name of the symbol containing this symbol.
	 */
	ContainerName string `json:"containerName,omitempty"`
}

/**
 * The parameters of a Workspace Symbol Request.
 */
type WorkplaceSymbolParams struct {
	/**
	 * A non-empty query string
	 */
	Query string `json:"query"`
}

/**
 * Contains additional diagnostic information about the context in which
 * a code action is run.
 */
type CodeActionContext struct {
	/**
	 * An array of diagnostics.
	 */
	Diagnostics []Diagnostic `json:"diagnostics"`
}

/**
 * Params for the CodeActionRequest
 */
type CodeActionParams struct {
	/**
	 * The document in which the command was invoked.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The range for which the command was invoked.
	 */
	Range Range `json:"range"`

	/**
	 * Context carrying additional information.
	 */
	Context CodeActionContext `json:"context"`
}

type CodeLensParams struct {
	/**
	 * The document to request code lens for.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`
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
	Range Range `json:"range"`

	/**
	 * The command this code lens represents.
	 */
	Command Command `json:"command,omitempty"`

	/**
	 * A data entry field that is preserved on a code lens item between
	 * a code lens and a code lens resolve request.
	 */
	Data interface{} `json:"data,omitempty"`
}

type DocumentFormattingParams struct {
	/**
	 * The document to format.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The format options.
	 */
	Options FormattingOptions `json:"options"`
}

/**
 * Value-object describing what options formatting should use.
 */
type FormattingOptions struct {
	/**
	 * Size of a tab in spaces.
	 */
	TabSize int32 `json:"tabSize"`

	/**
	 * Prefer spaces over tabs.
	 */
	InsertSpaces bool `json:"insertSpaces"`

	/**
	* Signature for further properites
	 */
	Key string `json:"key"`
}

type RenameParams struct {
	/**
	 * The document to format.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The position at which this request was sent.
	 */
	Position Position `json:"position"`

	/**
	 * The new name of the symbol. If the given name is not valid the
	 * request must return a [ResponseError](#ResponseError) with an
	 * appropriate message set.
	 */
	NewName string `json:"newName"`
}
