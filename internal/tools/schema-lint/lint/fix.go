// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"fmt"
	"os"
	"strings"
)

// Rename is an auto-applicable fix that renames a property by replacing every
// quoted occurrence of From with To in the file. Anchoring on the surrounding
// double quotes keeps the replacement scoped to the property name as it appears
// in the schema key, d.Get/d.Set calls, tfschema tags and cross-field
// references, without touching unrelated substrings.
type Rename struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// AppliedRename records a rename that Apply performed on a file.
type AppliedRename struct {
	File string
	From string
	To   string
}

// Apply rewrites files in place to remediate the auto-applicable findings —
// property renames — and returns the renames it performed. Findings are
// processed in their existing order; if two findings rename the same name the
// first wins, and re-running the linter converges on any follow-on findings.
func Apply(findings []Finding) ([]AppliedRename, error) {
	type edit struct{ from, to string }
	byFile := map[string][]edit{}
	var order []string
	for _, f := range findings {
		if f.Rename == nil {
			continue
		}
		if _, seen := byFile[f.File]; !seen {
			order = append(order, f.File)
		}
		byFile[f.File] = append(byFile[f.File], edit{f.Rename.From, f.Rename.To})
	}

	var applied []AppliedRename
	for _, file := range order {
		content, err := os.ReadFile(file)
		if err != nil {
			return applied, fmt.Errorf("reading %s: %w", file, err)
		}
		info, err := os.Stat(file)
		if err != nil {
			return applied, fmt.Errorf("stat %s: %w", file, err)
		}

		text := string(content)
		done := map[string]bool{}
		for _, e := range byFile[file] {
			if e.from == "" || e.to == "" || e.from == e.to || done[e.from] {
				continue
			}
			quotedFrom := `"` + e.from + `"`
			if !strings.Contains(text, quotedFrom) {
				continue
			}
			text = strings.ReplaceAll(text, quotedFrom, `"`+e.to+`"`)
			done[e.from] = true
			applied = append(applied, AppliedRename{File: file, From: e.from, To: e.to})
		}

		if text != string(content) {
			if err := os.WriteFile(file, []byte(text), info.Mode().Perm()); err != nil {
				return applied, fmt.Errorf("writing %s: %w", file, err)
			}
		}
	}
	return applied, nil
}

// Unfixed returns the findings that Apply cannot remediate automatically, i.e.
// everything except property renames.
func Unfixed(findings []Finding) []Finding {
	out := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if f.Rename == nil {
			out = append(out, f)
		}
	}
	return out
}
