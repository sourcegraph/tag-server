package main

import (
	"log"
	"os"

	"github.com/sourcegraph/tag-server/ctags"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func run() error {
	p, err := ctags.Parse2(nil)
	if err != nil {
		return err
	}
	tags := p.Tags()

	for _, tag := range tags {
		log.Printf("    %+v", tag)
	}

	return nil
}
