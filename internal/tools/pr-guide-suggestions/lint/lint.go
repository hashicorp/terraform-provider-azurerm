// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package lint parses Go source without type information, builds the schema
// tree, runs the rule checks, and optionally restricts findings to lines changed
// since a base git ref.
package lint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

// Finding is a single rule violation with enough context to report it.
type Finding struct {
	RuleID   string         `json:"ruleId"`
	RuleName string         `json:"ruleName"`
	Severity rules.Severity `json:"severity"`
	Resource string         `json:"resource,omitempty"`
	Kind     string         `json:"kind,omitempty"`
	Path     string         `json:"path,omitempty"`
	File     string         `json:"file"`
	Line     int            `json:"line"`
	Column   int            `json:"column"`
	Message  string         `json:"message"`
	Fix      string         `json:"fix,omitempty"`
	// Rename, when non-nil, is an auto-applicable fix that -w performs by
	// replacing every quoted occurrence of the old name with the new one.
	Rename *Rename `json:"rename,omitempty"`
}

// Options configures a lint run.
type Options struct {
	// Only, when non-empty, restricts the run to these rule IDs.
	Only map[string]bool
	// Disable lists rule IDs to skip (takes precedence over Only).
	Disable map[string]bool
	// Changes, when non-nil, restricts findings to added/changed lines.
	Changes *Changes
}

// Run lints the given targets (files or directories) and returns findings.
func Run(targets []string, opts Options) ([]Finding, error) {
	files, err := gatherFiles(targets)
	if err != nil {
		return nil, err
	}

	active := activeRules(opts)

	fset := token.NewFileSet()
	byDir := make(map[string][]*ast.File)
	for _, f := range files {
		parsed, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
		dir := filepath.Dir(f)
		byDir[dir] = append(byDir[dir], parsed)
	}

	var findings []Finding
	for _, dir := range sortedKeys(byDir) {
		res := schematree.Build(fset, byDir[dir])
		for _, r := range active {
			r.Check(res, func(n *schematree.Node, message string, fix *rules.Fix) {
				if opts.Changes != nil && !changedInSpan(opts.Changes, fset, n) {
					return
				}
				pos := fset.Position(n.Key.Pos())
				f := Finding{
					RuleID:   r.ID,
					RuleName: r.Name,
					Severity: r.Severity,
					Resource: n.Resource,
					Kind:     n.Kind,
					Path:     n.Path,
					File:     pos.Filename,
					Line:     pos.Line,
					Column:   pos.Column,
					Message:  message,
				}
				if fix != nil {
					f.Fix = fix.Suggestion
					if fix.Rename != "" {
						f.Rename = &Rename{From: n.Name, To: fix.Rename}
					}
				}
				findings = append(findings, f)
			})
		}
	}

	sortFindings(findings)
	return dedupeFindings(findings), nil
}

// dedupeFindings collapses findings that share a rule and exact source position.
// A property's map key has a unique position, so two findings with the same
// (rule, file, line, column) are the same violation reached via different walk
// paths (for example a schema helper linted through several callers); reporting
// it once, at its source line, is what a developer fixes.
func dedupeFindings(in []Finding) []Finding {
	seen := make(map[string]struct{}, len(in))
	out := make([]Finding, 0, len(in))
	for _, f := range in {
		key := fmt.Sprintf("%s\x00%s\x00%d\x00%d", f.RuleID, f.File, f.Line, f.Column)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, f)
	}
	return out
}

// activeRules returns the rules enabled by opts.
func activeRules(opts Options) []*rules.Rule {
	var out []*rules.Rule
	for _, r := range rules.Registry {
		if opts.Disable[r.ID] {
			continue
		}
		if len(opts.Only) > 0 && !opts.Only[r.ID] {
			continue
		}
		out = append(out, r)
	}
	return out
}

// gatherFiles expands targets into a de-duplicated list of non-test .go files.
func gatherFiles(targets []string) ([]string, error) {
	seen := map[string]bool{}
	var files []string
	add := func(p string) {
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		if !seen[abs] {
			seen[abs] = true
			files = append(files, abs)
		}
	}

	for _, t := range targets {
		info, err := os.Stat(t)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			add(t)
			continue
		}
		err = filepath.WalkDir(t, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			add(path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(files)
	return files, nil
}

// changedInSpan reports whether any line of the property's definition (its key
// through the end of its value) was changed. Using the whole span rather than
// only the anchored key line keeps a finding when a change modifies a property's
// body — for example removing its ValidateFunc — even though the key line itself
// is unchanged.
func changedInSpan(changes *Changes, fset *token.FileSet, n *schematree.Node) bool {
	start := fset.Position(n.Key.Pos()).Line
	end := fset.Position(n.Value.End()).Line
	return changes.AddedInRange(fset.Position(n.Key.Pos()).Filename, start, end)
}

func sortedKeys(m map[string][]*ast.File) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortFindings(f []Finding) {
	sort.SliceStable(f, func(i, j int) bool {
		if f[i].File != f[j].File {
			return f[i].File < f[j].File
		}
		if f[i].Line != f[j].Line {
			return f[i].Line < f[j].Line
		}
		if f[i].Column != f[j].Column {
			return f[i].Column < f[j].Column
		}
		return f[i].RuleID < f[j].RuleID
	})
}
