package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strings"

	gorpc "github.com/gorilla/rpc"
	gojson "github.com/gorilla/rpc/json"
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
	Std bool `long:"std" description:"if true, listen on stdin and use the protocol that VSCode LSP plugins use; if false, run HTTP server"`
}

func (c *ServeCmd) Execute(args []string) error {
	if c.Std {
		// Write to tmp file when run on stdin/out (assume you are
		// being run by a VSCode plugin)
		out, err := os.OpenFile("/tmp/lsp.out", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		log.SetOutput(out)

		if err := rpc.RegisterName("lsp", ctags.NewLangService()); err != nil {
			return err
		}
		stdPipe := ReadWriteCloser{os.Stdin, os.Stdout}
		log.Printf("listening on stdin")
		for {
			rpc.ServeCodec(newStdCodec(stdPipe))
		}
		return nil
	} else {
		s := gorpc.NewServer()
		s.RegisterCodec(&codec{gojson.NewCodec()}, "application/json")
		if err := s.RegisterService(lsp.NewHTTPLangService(ctags.NewLangService()), "lsp"); err != nil {
			return err
		}
		http.Handle("/", s)

		log.Printf("listening on :9090")
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
	var rErr, wErr error
	rErr = r.ReadCloser.Close()
	wErr = r.WriteCloser.Close()
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

type bufReadCloser struct {
	*bufio.Reader
}

func (r bufReadCloser) Close() error { return nil }

// stdCodec is a net/rpc Codec for using when interacting with VSCode over stdin/out
type stdCodec struct {
	conn      io.ReadWriteCloser
	bufReader *bufio.Reader
	bufWriter *bufio.Writer

	dec *json.Decoder

	req *jsonRPCReq // temporary working state
}

var _ rpc.ServerCodec = (*stdCodec)(nil)

func newStdCodec(conn io.ReadWriteCloser) *stdCodec {
	bufReader := bufio.NewReader(conn)
	bufWriter := bufio.NewWriter(conn)
	return &stdCodec{
		conn:      conn,
		bufReader: bufReader,
		bufWriter: bufWriter,
		dec:       json.NewDecoder(bufReader),
	}
}

type jsonRPCReq struct {
	Method string           `json:"method"`
	Id     uint64           `json:"id"`
	Params *json.RawMessage `json:"params"`
}

type jsonRPCResp struct {
	Id     uint64      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func (c *stdCodec) ReadRequestHeader(rpcReq *rpc.Request) error {
	for {
		l, err := c.bufReader.ReadString('\n')
		if err != nil {
			return err
		}
		if l == "\r\n" || l == "\n" {
			break
		}
	}

	var req jsonRPCReq
	if err := c.dec.Decode(&req); err != nil {
		return err
	}
	rpcReq.ServiceMethod = req.Method
	rpcReq.Seq = req.Id
	c.req = &req

	switch rpcReq.ServiceMethod {
	case "textDocument/hover":
		rpcReq.ServiceMethod = "lsp.HoverRequest"
	default:
		rpcReq.ServiceMethod = fmt.Sprintf("lsp.%s", capitalize(rpcReq.ServiceMethod))
	}

	return nil
}

func (c *stdCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if c.req.Params == nil {
		return errors.New("jsonrpc: request body missing params")
	}
	return json.Unmarshal(*c.req.Params, x)
}

// WriteResponse must be safe for concurrent use by multiple goroutines.
func (c *stdCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	resp := jsonRPCResp{Id: r.Seq}
	if r.Error == "" {
		resp.Result = x
	} else {
		resp.Error = r.Error
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.bufWriter, "Content-Length: %d\r\nContent-Type: application/vscode-jsonrpc; charset=utf8\r\n\r\n", len(body))
	if _, err := c.bufWriter.Write(body); err != nil {
		return err
	}
	return c.bufWriter.Flush()
}

func (c *stdCodec) Close() error {
	return c.conn.Close()
}

// codec is a gorilla/rpc Codec for use in serving JSON-RPC requests
// over HTTP
type codec struct {
	*gojson.Codec
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
