package ctags

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	File          string
	DefLinePrefix string
	Name          string

	// Extension fields
	Access         string // "private", "public"
	FileScope      string // ?
	Inheritance    string // ?
	Kind           string // "class"
	Language       string // "Java"
	Implementation string // ?
	Line           int    // 23
	Scope          string // "enum:gl::foobar"
	Signature      string // "(rtclass,objtype,obj,hr)"
	Type           string // ?
}

type TagsParser struct {
	// input
	config *Config

	// output
	tags      []Tag
	langFiles map[string][]string

	// temporary state
	curFile string
}

func NewParser2() (*TagsParser, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	return &TagsParser{
		langFiles: make(map[string][]string),
		config:    cfg,
	}, nil
}

func (p *TagsParser) Tags() []Tag {
	return p.tags
}

func (p *TagsParser) Parse(r *bufio.Reader) error {
	p.curFile = ""

	line, err := r.ReadString('\n')
	for ; err == nil; line, err = r.ReadString('\n') {
		if err := p.parseLine(strings.TrimRight(line, "\r\n")); err != nil {
			return err
		}
	}
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (p *TagsParser) parseLine(line string) error {
	if len(strings.TrimSpace(line)) == 0 || strings.HasPrefix(line, "!") {
		return nil
	}

	t1 := strings.Index(line, "\t")
	if t1 == -1 {
		return fmt.Errorf("expected tab-delimited line with at least 4 fields, but got %q", line)
	}
	name := line[0:t1]

	t2_ := strings.Index(line[t1+1:], "\t")
	if t2_ == -1 {
		return fmt.Errorf("expected tab-delimited line with at least 4 fields, but got %q", line)
	}
	t2 := t1 + 1 + t2_
	file := line[t1+1 : t2]

	t3_ := strings.LastIndex(line[t2+1:], `;"`)
	if t3_ == -1 {
		return fmt.Errorf(`expected find command to terminate with ';"', but got %q`, line)
	}
	t3 := t3_ + 2 + t2 + 1
	if len(line) <= t3 || line[t3] != '\t' {
		return fmt.Errorf(`expected tab immediately following ';"', but got %q, line: was %q`, line[t3:t3+1], line)
	}
	findCmd := line[t2+1 : t3]

	extFields_ := strings.Split(line[t3+1:], "\t")
	extFields := make(map[string]string)
	for _, extField := range extFields_ {
		s := strings.Index(extField, ":")
		key, val := extField[0:s], extField[s+1:]
		extFields[key] = val
	}
	lineno, err := strconv.Atoi(extFields["line"])
	if err != nil {
		return fmt.Errorf("could not parse line number, line was %q", line)
	}

	p.tags = append(p.tags, Tag{
		Name:          name,
		File:          file,
		DefLinePrefix: findCmdToDefLinePrefix(findCmd),
		Access:        extFields["access"],
		// FileScope:      string,
		// Inheritance:    string,
		Kind:     extFields["kind"],
		Language: extFields["language"],
		// Implementation: string,
		Line:      lineno,
		Scope:     extFields["scope"],
		Signature: extFields["signature"],
		Type:      extFields["typeref"],
	})
	return nil
}

func findCmdToDefLinePrefix(findCmd string) string {
	def := strings.TrimSuffix(strings.TrimPrefix(findCmd, `/^`), `/;"`)
	if strings.HasSuffix(def, "$") {
		def = strings.TrimSuffix(def, "$")
	}
	return def
}

func Parse2(files []string) (*TagsParser, error) {
	const tagsFilename = "tags"
	args := []string{"-f", tagsFilename, "--fields=*", "--excmd=pattern"}
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
	p, err := NewParser2()
	if err != nil {
		return nil, err
	}
	if err := p.Parse(r); err != nil {
		return nil, err
	}
	return p, nil
}
