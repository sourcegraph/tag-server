package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"sourcegraph.com/sourcegraph/sourcegraph/pkg/jsonrpc2"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/lsp"
)

type Config struct {
	Mode    string
	Addr    string
	Logfile string
}

func Serve(c Config) error {
	if c.Logfile != "" {
		f, err := os.Create(c.Logfile)
		if err != nil {
			return err
		}
		defer f.Close()
		log.SetOutput(io.MultiWriter(os.Stderr, f))
	}

	h := &jsonrpc2.LoggingHandler{Handler{}}

	switch c.Mode {
	case "tcp":
		lis, err := net.Listen("tcp", c.Addr)
		if err != nil {
			return err
		}
		defer lis.Close()
		log.Println("listening on", c.Addr)
		return jsonrpc2.Serve(lis, h)

	case "stdio":
		log.Println("reading on stdin, writing on stdout")
		jsonrpc2.NewServerConn(os.Stdin, os.Stdout, h)
		select {}

	default:
		return fmt.Errorf("invalid mode %q", c.Mode)
	}
}

type Handler struct{}

func (Handler) Handle(req *jsonrpc2.Request) (resp *jsonrpc2.Response) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("!!! PANIC recovered in Handle: %v", r)
		}
	}()

	if !req.Notification {
		resp = &jsonrpc2.Response{ID: req.ID}
	}

	switch req.Method {
	case "initialize":
		var res lsp.InitializeResult
		Server.Initialize(&lsp.InitializeParams{}, &res)
		resp.SetResult(res)

	case "shutdown":
		// Result is undefined, per
		// https://github.com/Microsoft/language-server-protocol/blob/master/protocol.md#shutdown-request.
		resp.SetResult(true)

	case "textDocument/hover":
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			resp.Error = &jsonrpc2.Error{Code: 123, Message: "error!"}
			return
		}

		pos := params.Position
		resp.SetResult(lsp.Hover{
			Contents: []lsp.MarkedString{{Language: "markdown", Value: "Hello CTags!"}},
			Range: lsp.Range{
				Start: lsp.Position{Line: pos.Line, Character: pos.Character - 3},
				End:   lsp.Position{Line: pos.Line, Character: pos.Character + 3},
			},
		})

	case "textDocument/documentSymbol":
		var params lsp.DocumentSymbolParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			resp.Error = &jsonrpc2.Error{Code: 123, Message: "error!"}
			return
		}

		var res []lsp.SymbolInformation
		Server.DocumentSymbols(&params, &res)
		resp.SetResult(res)

	case "textDocument/definition":
		var params lsp.TextDocumentPositionParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			resp.Error = &jsonrpc2.Error{Code: 123, Message: "error!"}
			return
		}

		var res []lsp.Location
		Server.GoToDefinition(&params, &res)
		resp.SetResult(res)

	case "textDocument/references":
		var params lsp.ReferenceParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			resp.Error = &jsonrpc2.Error{Code: 123, Message: "error!"}
			return
		}
		var res []lsp.Location
		Server.References(&params, &res)
		resp.SetResult(res)

	default:
		log.Printf("! Unrecognized RPC call: %s", req.Method)
	}

	return
}
