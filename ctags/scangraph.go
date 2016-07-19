package ctags

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"sourcegraph.com/sourcegraph/srclib/unit"
)

func Graph(files []string) (*Output, error) {
	p, err := Parse(files)
	if err != nil {
		return nil, err
	}
	return &Output{
		Defs: p.Defs(),
		Refs: p.Refs(),
	}, nil
}

func Scan() ([]*unit.SourceUnit, error) {
	p, err := Parse(nil)
	if err != nil {
		return nil, err
	}
	return p.Units(), nil
}

var ignoreFiles = []string{".srclib-cache", "node_modules", "vendor"}

func Parse(files []string) (*ETagsParser, error) {
	const tagsFilename = "tags"
	args := []string{"-e", "-f", tagsFilename}
	if len(files) == 0 {
		args = append(args, "-R")
	} else {
		args = append(args, files...)
	}
	excludeArgs := make([]string, 0, len(ignoreFiles))
	for _, ignoreFile := range ignoreFiles {
		excludeArgs = append(excludeArgs, fmt.Sprintf("--exclude=%s", ignoreFile))
	}
	args = append(args, excludeArgs...)

	log.Printf("...running ctags with args %v", args)
	ctagsStartTime := time.Now()
	if err := exec.Command("ctags", args...).Run(); err != nil {
		return nil, err
	}
	log.Printf("...done running ctags (duration: %v)", time.Since(ctagsStartTime))

	tagsFile, err := os.Open(tagsFilename)
	if err != nil {
		return nil, err
	}
	defer tagsFile.Close()

	r := bufio.NewReader(tagsFile)
	p, err := NewParser()
	if err != nil {
		return nil, err
	}
	if err := p.Parse(r); err != nil {
		return nil, err
	}
	return p, nil
}
