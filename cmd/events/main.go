package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/sourcegraph/tag-server/ctags"
)

var flagParser = flags.NewNamedParser("srclib-ctags", flags.Default)

func init() {
	_, err := flagParser.AddCommand("events",
		"output events",
		"output stream of events associated with HEAD commit",
		&eventsCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

var eventsCmd = EventsCmd{}

type EventsCmd struct{}

func main() {
	log.SetFlags(0)
	if _, err := flagParser.Parse(); err != nil {
		os.Exit(1)
	}
}

type Line struct {
	Num  int
	Text string
}

type HunkDiff struct {
	Filename string

	OldStart int
	OldEnd   int
	Old      []Line

	NewStart int
	NewEnd   int
	New      []Line
}

var fileHeaderRx = regexp.MustCompile(`diff \-\-git a\/([^\s]+) b\/(?:[^\s]+)`)
var hunkHeaderRx = regexp.MustCompile(`\@\@ \-([0-9]+),([0-9]+) \+([0-9]+),([0-9]+) \@\@`)
var typescriptRx = regexp.MustCompile(`<([A-Z]\w+).`)
var functionRx = regexp.MustCompile(`(?:([A-Za-z0-9]+)*\()`)
var branchRx = regexp.MustCompile(`HEAD branch: ([A-Za-z0-9]+)\n`)
var remoteRx = regexp.MustCompile(`Fetch\s*URL:\s*([A-Za-z0-9\.@:/-]+)\n`)

var ignore = map[string]bool{
	// go builtins and other ignore strings
	"append":     true,
	"cap":        true,
	"close":      true,
	"copy":       true,
	"delete":     true,
	"image":      true,
	"len":        true,
	"make":       true,
	"new":        true,
	"print":      true,
	"panic":      true,
	"println":    true,
	"real":       true,
	"recover":    true,
	"bool":       true,
	"byte":       true,
	"complex128": true,
	"complex64":  true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"int8":       true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uint8":      true,
	"uintptr":    true,
	"func":       true,
	"TODO":       true,
}

func generateURL(repository string, commitHash string) string {
	repository = strings.Replace(repository, "sourcegraph.com", "github.com", -1)
	return fmt.Sprintf("https://www.%s/commit/%s", repository, commitHash)
}

func (c *EventsCmd) Execute(args []string) error {
	var remoteURL, commitURL, commitHash, branch string
	{
		b, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			return err
		}
		commitHash = strings.TrimSpace(string(b))

		b, err = exec.Command("git", "remote", "show", "origin").Output()
		if err != nil {
			return err
		}
		branch = branchRx.FindStringSubmatch(string(b))[1]
		remoteURL = strings.Replace(
			strings.Replace(branch, "git@", "", 1), ":", "/", 1,
		)
		commitURL = generateURL(remoteURL, commitHash)
	}

	var hunkDiffs []*HunkDiff
	{
		// TODO(beyang): this introduces an off-by-one error, but we use unified=1 because it makes the hunk header regex simpler
		b, err := exec.Command("git", "show", "--unified=1").Output()
		if err != nil {
			return err
		}
		lines := strings.Split(string(b), "\n")
		oldline, newline := -1, -1 // keep track of current lines in new and old
		filename := ""
		for _, line := range lines {
			// detect file header
			if fileHeader := fileHeaderRx.FindStringSubmatch(line); len(fileHeader) == 2 {
				filename = fileHeader[1]
				// fileDiffs = append(fileDiffs, &FileDiff{Filename: fileHeader[1]})
				oldline, newline = -1, -1
				continue
			}
			// ignore if first file not yet found
			if filename == "" {
				continue
			}
			// ignore metadata lines
			if strings.HasPrefix(line, "diff --git") || strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
				continue
			}

			if hunkHeader := hunkHeaderRx.FindStringSubmatch(line); len(hunkHeader) == 5 {
				oldstart, _ := strconv.Atoi(hunkHeader[1])
				oldoff, _ := strconv.Atoi(hunkHeader[2])
				oldend := oldstart + oldoff - 1
				newstart, _ := strconv.Atoi(hunkHeader[3])
				newoff, _ := strconv.Atoi(hunkHeader[4])
				newend := newstart + newoff - 1
				oldline, newline = oldstart, oldend
				hunkDiffs = append(hunkDiffs, &HunkDiff{Filename: filename, OldStart: oldstart, OldEnd: oldend, NewStart: newstart, NewEnd: newend})
				continue
			}
			// ignore if first hunk not yet found
			if len(hunkDiffs) == 0 {
				continue
			}
			if strings.HasPrefix(line, "+") {
				hd := hunkDiffs[len(hunkDiffs)-1]
				hd.New = append(hd.New, Line{Num: newline, Text: line})
				newline++
			} else if strings.HasPrefix(line, "-") {
				hd := hunkDiffs[len(hunkDiffs)-1]
				hd.Old = append(hd.Old, Line{Num: oldline, Text: line})
				oldline++
			} else {
				oldline++
				newline++
			}
		}
	}

	var events []*EvtUpdate
	{ // definition modification events
		// TODO(beyang): include authorship information for each def
		files := make([]string, 0, len(hunkDiffs))
		for _, hd := range hunkDiffs {
			if len(files) == 0 || files[len(files)-1] != hd.Filename {
				files = append(files, hd.Filename)
			}
		}
		hunkDiffM := make(map[string][]*HunkDiff)
		for _, hd := range hunkDiffs {
			hunkDiffM[hd.Filename] = append(hunkDiffM[hd.Filename], hd)
		}

		p, err := ctags.Parse2(files)
		if err != nil {
			return err
		}

		tags := p.Tags()
		sort.Sort(tagSorter{tags})
		var changedTags []*ctags.Tag
		for i, _ := range tags {
			endline := math.MaxInt64
			if i+1 < len(tags) {
				endline = tags[i+1].Line - 1
			}
			for _, hd := range hunkDiffM[tags[i].File] {
				if !(hd.NewStart > endline || hd.NewEnd < tags[i].Line) {
					// tag overlaps with diff
					changedTags = append(changedTags, &tags[i])
					break
				}
			}
		}
		for _, tag := range changedTags {
			events = append(events, &EvtUpdate{
				Hashes: []string{tag.Name, tag.File},
				Users:  nil,
				Event: &Evt{
					ID:    fmt.Sprintf("modified:%s:%s:%s", tag.Name, tag.File, commitURL),
					Title: fmt.Sprintf("%s %s%s was modified", tag.Kind, tag.Name, tag.Signature),
					Body:  fmt.Sprintf("%s %s%s in %s was modified in commit", tag.Kind, tag.Name, tag.Signature, tag.File),
					URL:   commitURL,
					Type:  "modified",
					// TODO(beyang): time
				},
			})
		}
	}
	{ // reference events
		for _, hd := range hunkDiffs {
			for _, newLine := range hd.New {
				for _, match := range functionRx.FindAllStringSubmatch(newLine.Text, -1) {
					// temporary fix for bad regex, gr... regexes...
					if len(match[1]) > 0 && !ignore[match[1]] {
						events = append(events, &EvtUpdate{
							Hashes: []string{match[1], hd.Filename},
							Users:  nil,
							Event: &Evt{
								ID:    fmt.Sprintf("referenced:%s:%s:%s", match[1], hd.Filename, commitURL),
								Title: fmt.Sprintf("function %s was referenced", match[1]),
								Body:  fmt.Sprintf("function %s was referenced in file %s in commit %s on branch %s", match[1], hd.Filename, commitHash, branch),
								URL:   commitURL,
								Type:  "referenced",
							},
						})
					}
				}
				for _, match := range typescriptRx.FindStringSubmatch(newLine.Text) {
					if len(match) > 0 && !ignore[match] {
						events = append(events, &EvtUpdate{
							Hashes: []string{match, hd.Filename},
							Users:  nil,
							Event: &Evt{
								ID:    fmt.Sprintf("referenced(react):%s:%s:%s", match, hd.Filename, commitURL),
								Title: fmt.Sprintf("React component %s was used", match),
								Body:  fmt.Sprintf("React component %s was used in file %s in commit %s on branch %s", match, hd.Filename, commitHash, branch),
								URL:   commitURL,
								Type:  "referenced",
							},
						})
					}
				}
			}
		}
	}

	return json.NewEncoder(os.Stdout).Encode(EvtsPostOpts{Updates: events})
}

type tagSorter struct {
	tags []ctags.Tag
}

func (t tagSorter) Less(i, j int) bool {
	return t.tags[i].Line < t.tags[j].Line
}
func (t tagSorter) Swap(i, j int) {
	t.tags[i], t.tags[j] = t.tags[j], t.tags[i]
}
func (t tagSorter) Len() int {
	return len(t.tags)
}
