package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sourcegraph/tag-server/server"
)

var (
	mode    = flag.String("mode", "stdio", "communication mode (stdio|tcp)")
	addr    = flag.String("addr", ":2088", "server listen address (tcp)")
	logfile = flag.String("log", "/tmp/sample_server.log", "write log output to this file (and stderr)")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	if err := server.Serve(server.Config{
		Mode:    *mode,
		Addr:    *addr,
		Logfile: *logfile,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
