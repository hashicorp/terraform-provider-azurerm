// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

// Fixer fix document with diff
type Fixer struct {
	MDFile       string
	SchemaFile   string
	ResourceType string

	Diff []Checker // diff of exist in both md and code

	FixedContent string
}

func NewFixer(d *ResourceDiff) *Fixer {
	f := &Fixer{
		MDFile:       d.MDFile,
		SchemaFile:   d.SchemaFile,
		ResourceType: d.tf.ResourceType,
		Diff:         d.Diff,
	}
	return f
}

// param rt: resource type
// param lines: the lines of markdown file
func tryFixTimeouts(rt string, lines []string, diffs []TimeoutDiffItem) []string {
	var suf []string
	addSuf := func(line string) {
		suf = append(suf, line)
	}
	if diffs[0].Type == TimeoutMissed {
		// no such timeout block, add to the end of lines
		addSuf("## Timeouts")
		addSuf("")
		addSuf("The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:")
		addSuf("")
		diffs = diffs[1:]
	}
	// find timeout block
	var toLine, importIdx int
	for idx, line := range lines {
		if line == "## Timeouts" {
			toLine = idx + 4
			for toLine < len(lines) && lines[toLine] != "" {
				toLine++ // insert to an empty line
			}
		}
		if strings.HasPrefix(line, "## Import") {
			importIdx = idx
		}
	}
	rt = util.NormalizeResourceName(rt)

	for _, diff := range diffs {
		if diff.Line == 0 {
			// append line
			gen := diff.GenLine(rt)
			if len(suf) > 0 {
				addSuf(gen)
			} else {
				lines = append(lines[:toLine+1], lines[toLine:]...)
				lines[toLine] = gen
			}
		} else {
			lines[diff.Line-1] = diff.FixLine(lines[diff.Line-1])
		}
	}
	if len(suf) > 0 {
		addSuf("")
		// insert before import
		if importIdx > 0 {
			end := make([]string, len(lines)-importIdx)
			copy(end, lines[importIdx:])
			lines = append(lines[:importIdx], suf...)
			lines = append(lines, end...)
		} else {
			lines = append(lines, suf...)
		}
	}
	return lines
}

func (f *Fixer) TryFix() (err error) {
	// read file as bytes
	if len(f.Diff) == 0 {
		return
	}
	if d, ok := f.Diff[0].(diffWithMessage); ok {
		log.Printf("%s: %s", f.ResourceType, d.msg)
		return
	}
	content, err := os.ReadFile(f.MDFile)
	if err != nil {
		log.Printf("open %s: %v", f.MDFile, err)
		return err
	}

	lines := strings.Split(string(content), "\n")
	for idx, item := range f.Diff {
		_ = idx
		// fix timeout first!
		if to, ok := item.(timeoutDiff); ok {
			lines = tryFixTimeouts(f.ResourceType, lines, to.TimeoutDiff)
			continue
		}

		// mdField is nil for no document exists or page title mismatch
		if item.ShouldSkip() {
			continue
		}

		lineIdx := item.Line()
		line := lines[lineIdx]

		if line, err = item.Fix(line); err != nil {
			return err
		}

		if suf := strings.TrimSuffix(line, " "); suf != "" {
			if ch := suf[len(suf)-1]; ch != '.' && ch != '?' {
				line = suf + "."
			}
		}

		lines[lineIdx] = line
	}
	f.FixedContent = strings.Join(lines, "\n")
	return nil
}

func (f *Fixer) WriteBack() (err error) {
	if len(f.Diff) == 0 {
		return
	}
	if f.FixedContent == "" {
		log.Printf("%s no content to write back, skip", f.MDFile)
		return
	}
	fd, err := os.OpenFile(f.MDFile, os.O_TRUNC|os.O_RDWR, 066)
	if err != nil {
		log.Printf("open %s: %v", f.MDFile, err)
		return err
	}
	defer func() {
		_ = fd.Sync()
		_ = fd.Close()
	}()
	_, err = fd.WriteString(f.FixedContent)
	if err != nil {
		log.Printf("write %s back: %v", f.MDFile, err)
		return err
	}
	return nil
}
