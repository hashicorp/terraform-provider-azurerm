// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// WriteText renders findings as human-readable lines relative to root. A
// suggested remediation is printed for any finding that carries one.
func WriteText(w io.Writer, findings []Finding, root string) {
	for _, f := range findings {
		loc := f.File
		if rel, err := filepath.Rel(root, f.File); err == nil && !strings.HasPrefix(rel, "..") {
			loc = rel
		}
		msg := f.Message
		if f.Resource != "" {
			msg = f.Resource + ": " + msg
		}
		fmt.Fprintf(w, "%s:%d:%d: %s [%s] %s\n", loc, f.Line, f.Column, f.Severity, f.RuleID, msg)
		if f.Fix != "" {
			fmt.Fprintf(w, "    → fix: %s\n", f.Fix)
		}
	}

	errors, warnings := counts(findings)
	fmt.Fprintf(w, "\n%d finding(s): %d error(s), %d warning(s)\n", len(findings), errors, warnings)
}

// WriteJSON renders findings as a JSON array.
func WriteJSON(w io.Writer, findings []Finding) error {
	if findings == nil {
		findings = []Finding{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(findings)
}

// WriteApplied reports the renames Apply performed, relative to root.
func WriteApplied(w io.Writer, applied []AppliedRename, root string) {
	for _, a := range applied {
		loc := a.File
		if rel, err := filepath.Rel(root, a.File); err == nil && !strings.HasPrefix(rel, "..") {
			loc = rel
		}
		fmt.Fprintf(w, "fixed %s: renamed %q → %q\n", loc, a.From, a.To)
	}
	if len(applied) > 0 {
		fmt.Fprintf(w, "\napplied %d fix(es)\n", len(applied))
	}
}

// HasErrors reports whether any finding is error severity.
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == "error" {
			return true
		}
	}
	return false
}

func counts(findings []Finding) (errors, warnings int) {
	for _, f := range findings {
		if f.Severity == "error" {
			errors++
		} else {
			warnings++
		}
	}
	return errors, warnings
}
