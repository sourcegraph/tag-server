package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	if err := s.RegisterService(&LangSvc{}, ""); err != nil {
		return err
	}
	http.Handle("/", s)

	fmt.Println("listening on :9090")
	http.ListenAndServe(":9090", nil)
	return nil
}

type LangSvc struct{}

type DoArgs struct{}
type DoReply struct{}

func (s *LangSvc) Do(req *http.Request, args *DoArgs, reply *DoReply) error {
	fmt.Println("HELLO")
	return nil
}
