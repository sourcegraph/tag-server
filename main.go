package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/sourcegraph/tag-server/ctags"
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
		"srclib graph phase",
		"srclib graph phase",
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
	os.Stdin.Close() // ignore input source unit

	out, err := ctags.Graph(c.Files)
	if err != nil {
		fmt.Printf("failed due to error: %s\n", err)
		os.Exit(1)
	}
	return json.NewEncoder(os.Stdout).Encode(out)
}

/*
 * Scan
 */
func init() {
	_, err := flagParser.AddCommand("scan",
		"srclib scan phase",
		"srclib scan phase",
		&scanCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type ScanCmd struct{}

var scanCmd ScanCmd

func (c *ScanCmd) Execute(args []string) error {
	units, err := ctags.Scan()
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(units)
}

/*
 * Depresolve
 */
func init() {
	_, err := flagParser.AddCommand("depresolve",
		"srclib depresolve phase",
		"srclib depresolve phase",
		&depresolveCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type DepresolveCmd struct{}

var depresolveCmd DepresolveCmd

func (c *DepresolveCmd) Execute(args []string) error {
	fmt.Println("[]")
	return nil
}
