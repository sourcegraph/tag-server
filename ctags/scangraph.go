package ctags

import (
	"bufio"
	"os"
	"os/exec"

	"sourcegraph.com/sourcegraph/srclib/unit"
)

func Graph(files []string) (*Output, error) {
	p, err := parse(files)
	if err != nil {
		return nil, err
	}
	return &Output{
		Defs: p.Defs(),
	}, nil
}

func Scan() ([]*unit.SourceUnit, error) {
	p, err := parse(nil)
	if err != nil {
		return nil, err
	}
	return p.Units(), nil
}

func parse(files []string) (*ETagsParser, error) {
	const tagsFilename = "tags"
	args := []string{"-e", "-f", tagsFilename}
	if len(files) == 0 {
		args = append(args, "-R")
	} else {
		args = append(args, files...)
	}
	if err := exec.Command("ctags", args...).Run(); err != nil {
		return nil, err
	}
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
