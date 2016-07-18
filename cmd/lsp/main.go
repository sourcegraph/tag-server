package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strings"

	gorpc "github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/jessevdk/go-flags"
	"github.com/sourcegraph/tag-server/ctags"
	"github.com/sourcegraph/tag-server/lsp"
)

func main() {
	cli := flags.NewNamedParser("lsp", flags.PrintErrors|flags.PassDoubleDash)
	cli.AddCommand("serve", "", "", &ServeCmd{})
	_, err := cli.Parse()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type ServeCmd struct {
	Std bool `long:"std" description:""`
}

func (c *ServeCmd) Execute(args []string) error {
	if c.Std {
		if err := rpc.Register(ctags.NewLangService()); err != nil {
			return err
		}
		stdPipe := ReadWriteCloser{os.Stdout, os.Stdin}
		fmt.Println("listening on stdin")
		rpc.ServeCodec(jsonrpc.NewServerCodec(stdPipe))
		return nil
	} else {
		s := gorpc.NewServer()
		s.RegisterCodec(&codec{json.NewCodec()}, "application/json")
		if err := s.RegisterService(lsp.NewHTTPLangService(ctags.NewLangService()), "lsp"); err != nil {
			return err
		}
		http.Handle("/", s)

		fmt.Println("listening on :9090")
		return http.ListenAndServe(":9090", nil)
	}
}

type ReadWriteCloser struct {
	io.ReadCloser
	io.WriteCloser
}

func (r ReadWriteCloser) Read(p []byte) (n int, err error) {
	return r.ReadCloser.Read(p)
}

func (r ReadWriteCloser) Write(p []byte) (n int, err error) {
	return r.WriteCloser.Write(p)
}

func (r ReadWriteCloser) Close() error {
	rErr := r.ReadCloser.Close()
	wErr := r.WriteCloser.Close()
	if rErr != nil {
		if wErr != nil {
			return fmt.Errorf("ReadCloser error: %s, WriteCloser error: %s", rErr, wErr)
		}
		return rErr
	}
	if wErr != nil {
		return wErr
	}
	return nil
}

type codec struct {
	*json.Codec
}

var _ gorpc.Codec = (*codec)(nil)

type codecReq struct {
	gorpc.CodecRequest
}

func (r *codecReq) Method() (string, error) {
	method, err := r.CodecRequest.Method()
	switch method {
	case "textDocument/hover":
		method = "lsp.HoverRequest"
	default:
		method = fmt.Sprintf("lsp.%s", capitalize(method))
	}
	return method, err
}

func (r *codecReq) ReadRequest(v interface{}) error {
	return r.CodecRequest.ReadRequest(v)
}

func (r *codecReq) WriteResponse(w http.ResponseWriter, v interface{}, err error) error {
	return r.CodecRequest.WriteResponse(w, v, err)
}

var _ gorpc.CodecRequest = (*codecReq)(nil)

func (c *codec) NewRequest(r *http.Request) gorpc.CodecRequest {
	return &codecReq{c.Codec.NewRequest(r)}
}

func capitalize(s string) string {
	if len(s) < 1 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
