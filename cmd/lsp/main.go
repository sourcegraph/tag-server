package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sourcegraph/tag-server/ctags"
	"github.com/sourcegraph/tag-server/lsp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	s := rpc.NewServer()
	s.RegisterCodec(&codec{json.NewCodec()}, "application/json")
	if err := s.RegisterService(lsp.NewHTTPLangService(ctags.NewLangService()), "lsp"); err != nil {
		return err
	}
	http.Handle("/", s)

	fmt.Println("listening on :9090")
	http.ListenAndServe(":9090", nil)
	return nil
}

type codec struct {
	*json.Codec
}

var _ rpc.Codec = (*codec)(nil)

type codecReq struct {
	rpc.CodecRequest
}

func (r *codecReq) Method() (string, error) {
	method, err := r.CodecRequest.Method()
	method = fmt.Sprintf("lsp.%s", capitalize(method))
	return method, err
}

func (r *codecReq) ReadRequest(v interface{}) error {
	return r.CodecRequest.ReadRequest(v)
}

func (r *codecReq) WriteResponse(w http.ResponseWriter, v interface{}, err error) error {
	return r.CodecRequest.WriteResponse(w, v, err)
}

var _ rpc.CodecRequest = (*codecReq)(nil)

func (c *codec) NewRequest(r *http.Request) rpc.CodecRequest {
	return &codecReq{c.Codec.NewRequest(r)}
}

func capitalize(s string) string {
	if len(s) < 1 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
