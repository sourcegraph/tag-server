package ctags

import "sourcegraph.com/sourcegraph/srclib/graph"

// Def is a definition in code.
type Def struct {
	// DefKey is the natural unique key for a def. It is stable
	// (subsequent runs of a grapher will emit the same defs with the same
	// DefKeys).
	graph.DefKey `protobuf:"bytes,1,opt,name=Key,embedded=Key" json:""`
	// Name of the definition. This need not be unique.
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name"`
	// Kind is the kind of thing this definition is. This is
	// language-specific. Possible values include "type", "func",
	// "var", etc.
	Kind     string `protobuf:"bytes,3,opt,name=Kind,proto3" json:"Kind,omitempty"`
	File     string `protobuf:"bytes,4,opt,name=File,proto3" json:"File"`
	DefStart uint32 `protobuf:"varint,5,opt,name=DefStart,proto3" json:"DefStart"`
	DefEnd   uint32 `protobuf:"varint,6,opt,name=DefEnd,proto3" json:"DefEnd"`
	// Exported is whether this def is part of a source unit's
	// public API. For example, in Java a "public" field is
	// Exported.
	Exported bool `protobuf:"varint,7,opt,name=Exported,proto3" json:"Exported,omitempty"`
	// Local is whether this def is local to a function or some
	// other inner scope. Local defs do *not* have module,
	// package, or file scope. For example, in Java a function's
	// args are Local, but fields with "private" scope are not
	// Local.
	Local bool `protobuf:"varint,8,opt,name=Local,proto3" json:"Local,omitempty"`
	// Test is whether this def is defined in test code (as opposed to main
	// code). For example, definitions in Go *_test.go files have Test = true.
	Test bool `protobuf:"varint,9,opt,name=Test,proto3" json:"Test,omitempty"`
	// Data contains additional language- and toolchain-specific information
	// about the def. Data is used to construct function signatures,
	// import/require statements, language-specific type descriptions, etc.
	Data interface{}
	// Docs are docstrings for this Def. This field is not set in the
	// Defs produced by graphers; they should emit docs in the
	// separate Docs field on the graph.Output struct.
	Docs []*graph.DefDoc `protobuf:"bytes,11,rep,name=Docs" json:"Docs,omitempty"`
	// TreePath is a structurally significant path descriptor for a def. For
	// many languages, it may be identical or similar to DefKey.Path.
	// However, it has the following constraints, which allow it to define a
	// def tree.
	//
	// A tree-path is a chain of '/'-delimited components. A component is either a
	// def name or a ghost component.
	// - A def name satifies the regex [^/-][^/]*
	// - A ghost component satisfies the regex -[^/]*
	// Any prefix of a tree-path that terminates in a def name must be a valid
	// tree-path for some def.
	// The following regex captures the children of a tree-path X: X(/-[^/]*)*(/[^/-][^/]*)
	TreePath string `protobuf:"bytes,17,opt,name=TreePath,proto3" json:"TreePath,omitempty"`
	Line     uint32 `protobuf:"varint,18,opt,name=Line,proto3" json:"Line"`
}

type Output struct {
	Defs []*Def
	Refs []*graph.Ref
	Docs []*graph.Doc
}
