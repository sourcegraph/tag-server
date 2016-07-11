package graph

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"sourcegraph.com/sourcegraph/srclib/graph"
	"sourcegraph.com/sqs/pbtypes"
)

func DefsForFiles(files []string) ([]*Def, error) {
	if len(files) == 0 {
		return nil, nil
	}

	args := []string{"-e", "-f", "tags"}
	args = append(args, files...)
	if err := exec.Command("ctags", args...).Run(); err != nil {
		return nil, err
	}
	tagsFile, err := os.Open("tags")
	if err != nil {
		return nil, err
	}
	defer tagsFile.Close()

	r := bufio.NewReader(tagsFile)
	tags, err := (&ETagsParser{}).Parse(r)
	if err != nil {
		return nil, err
	}

	return defsFromTags(tags), nil
}

func defsFromTags(tags []Tag) []*Def {
	defs := make([]*Def, len(tags))
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		formatData, _ := json.Marshal(defFormatDataFromTag(tag))
		defs[i] = &Def{
			DefKey: graph.DefKey{
				UnitType: "GoPackage", // dummy
				Unit:     "CTagsUnit", // dummy
				Path:     fmt.Sprintf("%s:%s", tag.File, tag.Name),
			},
			Name:     tag.Name,
			File:     tag.File,
			DefStart: uint32(tag.ByteOff),
			DefEnd:   uint32(tag.ByteOff + len(tag.Name)),
			Exported: true,
			Local:    false,
			Data:     pbtypes.RawMessage(formatData),
		}
	}
	return defs
}

// defFormatDataFromTag returns the display formatting data for a
// definition derived from the specified tag.
//
// Precondition: it assumes that tag.Name exists in tag.Def.
func defFormatDataFromTag(tag Tag) *DefFormatData {
	nameIdx := strings.Index(tag.Def, tag.Name)
	keyword := strings.TrimSpace(tag.Def[:nameIdx])
	typ := tag.Def[nameIdx+len(tag.Name):]
	sep := ""
	if len(typ) >= 1 && typ[0] == ':' {
		sep, typ = typ[0:1], strings.TrimSpace(typ[1:])
	}
	return &DefFormatData{
		Name:      tag.Name,
		Keyword:   keyword,
		Type:      typ,
		Kind:      keyword,
		Separator: sep,
	}
}

// This mirrors the format data (DefData) struct of Sourcegraph's
// basic def formatter. We don't depend directly on that because we
// should have no dependencies on Sourcegraph here.
type DefFormatData struct {
	Name      string
	Keyword   string
	Type      string
	Kind      string
	Separator string
}

type Tag struct {
	File    string
	Def     string
	Name    string
	Line    int
	ByteOff int
}

const (
	sepTag = "\x7f"
	sepPos = "\x01"
	sepCol = ","
)

type ETagsParser struct {
	curFile string
	tags    []Tag
}

func (p *ETagsParser) Parse(r *bufio.Reader) ([]Tag, error) {
	line, err := r.ReadString('\n')
	for ; err == nil; line, err = r.ReadString('\n') {
		if err := p.parseLine(strings.TrimSpace(line)); err != nil {
			return nil, err
		}
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	return p.tags, nil
}

func (p *ETagsParser) parseLine(line string) error {
	if len(line) == 0 || strings.HasPrefix(line, "!") {
		return nil
	}

	nameIdx := strings.Index(line, sepTag)
	if nameIdx < 0 {
		// File line
		cmps := strings.Split(line, ",")
		if len(cmps) != 2 {
			return fmt.Errorf("tags line parsing error: unrecognized format, line was %q", line)
		}
		if _, err := strconv.Atoi(cmps[1]); err != nil {
			return fmt.Errorf("tags line parsing error: %s, line was %q", err, line)
		}
		p.curFile = cmps[0]
		return nil
	}

	// Symbol line
	lineNoIdx_ := strings.Index(line[nameIdx:], sepPos)
	if lineNoIdx_ < 0 {
		return fmt.Errorf("tags line parsing error: could not find character %U, line was %q", sepPos, line)
	}
	lineNoIdx := nameIdx + lineNoIdx_

	colIdx_ := strings.Index(line[lineNoIdx:], sepCol)
	if colIdx_ < 0 {
		return fmt.Errorf("tags line parsing error: could not find character %q, line was %q", sepCol, line)
	}
	colIdx := lineNoIdx + colIdx_

	lineNo, err := strconv.Atoi(line[lineNoIdx+1 : colIdx])
	if err != nil {
		return fmt.Errorf("tags line parsing error: could not parse line number, line was %q", line)
	}
	colNo, err := strconv.Atoi(line[colIdx+1:])
	if err != nil {
		return fmt.Errorf("tags line parsing error: could not parse byte offset, line was %q", line)
	}

	p.tags = append(p.tags, Tag{
		File:    p.curFile,
		Def:     line[0:nameIdx],
		Name:    line[nameIdx+1 : lineNoIdx],
		Line:    lineNo,
		ByteOff: colNo,
	})
	return nil
}
