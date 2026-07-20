// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"go/ast"
	"go/token"
	"strings"
)

// nolintKeyword is the directive, modelled on the golangci-lint //nolint
// convention, that silences schema-lint findings for a property. A bare
// //nolint suppresses every rule; //nolint:SL001,SL002 suppresses only the
// listed rules (rule IDs are matched case-insensitively).
const nolintKeyword = "nolint"

// nolintDirective records which rules a //nolint comment silences.
type nolintDirective struct {
	all   bool
	rules map[string]bool
}

// nolintIndex resolves whether a rule is suppressed for the property whose map
// key is on a given source line. A directive applies to a property when it
// trails the key line (e.g. `"foo": { //nolint`) or stands on its own line
// directly above the key.
type nolintIndex struct {
	byLine map[string]map[int]*nolintDirective
}

// suppressed reports whether ruleID is silenced for the property key at
// file:line.
func (ix *nolintIndex) suppressed(file string, line int, ruleID string) bool {
	if ix == nil {
		return false
	}
	d := ix.byLine[file][line]
	if d == nil {
		return false
	}
	return d.all || d.rules[strings.ToUpper(ruleID)]
}

// buildNolintIndex scans the parsed files' comments for //nolint directives and
// maps each to the property key line it applies to. sources holds every file's
// lines (split on "\n") so a trailing directive can be told apart from one on
// its own line.
func buildNolintIndex(fset *token.FileSet, files []*ast.File, sources map[string][]string) *nolintIndex {
	ix := &nolintIndex{byLine: map[string]map[int]*nolintDirective{}}
	for _, f := range files {
		for _, group := range f.Comments {
			d, ok := groupDirective(group)
			if !ok {
				continue
			}
			start := fset.Position(group.List[0].Slash)
			target := start.Line
			if ownLine(sources[start.Filename], start.Line, start.Column) {
				// A standalone directive applies to the code line that follows
				// the whole comment group (so stacked directives all attach to
				// the property below them).
				target = fset.Position(group.End()).Line + 1
			}
			ix.add(start.Filename, target, d)
		}
	}
	return ix
}

// groupDirective merges every //nolint directive found in a comment group.
func groupDirective(group *ast.CommentGroup) (nolintDirective, bool) {
	var out nolintDirective
	found := false
	for _, c := range group.List {
		d, ok := parseNolint(c.Text)
		if !ok {
			continue
		}
		found = true
		out.all = out.all || d.all
		for r := range d.rules {
			if out.rules == nil {
				out.rules = map[string]bool{}
			}
			out.rules[r] = true
		}
	}
	return out, found
}

// parseNolint parses one comment's text into a directive. It accepts //nolint,
// //nolint:SL001,SL002 and a trailing " // explanation".
func parseNolint(text string) (nolintDirective, bool) {
	body, ok := commentBody(text)
	if !ok {
		return nolintDirective{}, false
	}
	field := body
	if i := strings.IndexAny(field, " \t"); i >= 0 {
		field = field[:i]
	}
	if field == nolintKeyword {
		return nolintDirective{all: true}, true
	}
	list, ok := strings.CutPrefix(field, nolintKeyword+":")
	if !ok {
		return nolintDirective{}, false
	}
	rules := map[string]bool{}
	for _, part := range strings.Split(list, ",") {
		if p := strings.ToUpper(strings.TrimSpace(part)); p != "" {
			rules[p] = true
		}
	}
	if len(rules) == 0 {
		return nolintDirective{all: true}, true
	}
	return nolintDirective{rules: rules}, true
}

// commentBody strips the comment markers and surrounding space from a raw
// comment.
func commentBody(text string) (string, bool) {
	switch {
	case strings.HasPrefix(text, "//"):
		text = text[2:]
	case strings.HasPrefix(text, "/*"):
		text = strings.TrimSuffix(strings.TrimPrefix(text, "/*"), "*/")
	default:
		return "", false
	}
	return strings.TrimSpace(text), true
}

// ownLine reports whether only whitespace precedes the given 1-based column,
// i.e. the comment stands on its own rather than trailing code.
func ownLine(lines []string, line, col int) bool {
	if line < 1 || line > len(lines) {
		return false
	}
	prefix := lines[line-1]
	if col-1 > len(prefix) {
		return true
	}
	return strings.TrimSpace(prefix[:col-1]) == ""
}

// add records a directive against a file line, merging with any existing one.
func (ix *nolintIndex) add(file string, line int, d nolintDirective) {
	byLine := ix.byLine[file]
	if byLine == nil {
		byLine = map[int]*nolintDirective{}
		ix.byLine[file] = byLine
	}
	cur := byLine[line]
	if cur == nil {
		dc := d
		byLine[line] = &dc
		return
	}
	cur.all = cur.all || d.all
	for r := range d.rules {
		if cur.rules == nil {
			cur.rules = map[string]bool{}
		}
		cur.rules[r] = true
	}
}
