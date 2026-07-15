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

// WriteText renders findings as human-readable lines relative to root. When
// showFix is true, a suggested remediation is printed for fixable findings.
func WriteText(w io.Writer, findings []Finding, root string, showFix bool) {
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
		if showFix && f.Fix != "" {
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
