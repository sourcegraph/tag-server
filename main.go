package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/beyang/srclib-ctags/graph"
	"github.com/jessevdk/go-flags"
)

var (
	flagParser = flags.NewNamedParser("srclib-ctags", flags.Default)
	cwd        = getCWD()
)

func init() {
	flagParser.LongDescription = "srclib-ctags performs static analysis by parsing ctags-style indexes."
}
func getCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

func main() {
	log.SetFlags(0)
	if _, err := flagParser.Parse(); err != nil {
		os.Exit(1)
	}
}

/*
 * Graph
 */
func init() {
	_, err := flagParser.AddCommand("graph",
		"graph a Go package",
		"Graph a Go package, producing all defs, refs, and docs.",
		&graphCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type GraphCmd struct {
	Files []string `short:"f" long:"files" description:"files to process; if empty, processes all files"`
}

var graphCmd GraphCmd

func (c *GraphCmd) Execute(args []string) error {
	defs, err := graph.DefsForFiles(c.Files)
	if err != nil {
		fmt.Printf("failed due to error: %s\n", err)
		os.Exit(1)
	}
	return json.NewEncoder(os.Stdout).Encode(defs)
}
