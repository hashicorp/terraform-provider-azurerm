// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package report renders lint findings in text or JSON form.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/rules"
)

// Format identifies an output format.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Summary aggregates finding counts by severity.
type Summary struct {
	Errors   int
	Warnings int
}

// Summarize counts findings by severity.
func Summarize(findings []rules.Finding) Summary {
	s := Summary{}
	for _, f := range findings {
		switch f.Severity {
		case rules.SeverityError:
			s.Errors++
		case rules.SeverityWarning:
			s.Warnings++
		}
	}

	return s
}

// Write renders findings to w in the requested format (defaults to text).
func Write(w io.Writer, findings []rules.Finding, format Format) error {
	if format == FormatJSON {
		return writeJSON(w, findings)
	}

	return writeText(w, findings)
}

func writeJSON(w io.Writer, findings []rules.Finding) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(findings)
}

func writeText(w io.Writer, findings []rules.Finding) error {
	byResource := make(map[string][]rules.Finding)
	order := make([]string, 0)
	for _, f := range findings {
		key := fmt.Sprintf("%s (%s)", f.ResourceType, f.Kind)
		if _, ok := byResource[key]; !ok {
			order = append(order, key)
		}
		byResource[key] = append(byResource[key], f)
	}
	sort.Strings(order)

	for _, key := range order {
		if _, err := fmt.Fprintln(w, key); err != nil {
			return err
		}
		for _, f := range byResource[key] {
			loc := f.Path
			if loc == "" {
				loc = "(resource)"
			}
			if _, err := fmt.Fprintf(w, "  %-7s  %-6s  %s: %s\n", f.Severity, f.RuleID, loc, f.Message); err != nil {
				return err
			}
			if f.FixSuggestion != "" {
				if _, err := fmt.Fprintf(w, "           \u2192 fix: %s\n", f.FixSuggestion); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	s := Summarize(findings)
	_, err := fmt.Fprintf(w, "%d problem(s) (%d error(s), %d warning(s))\n", len(findings), s.Errors, s.Warnings)

	return err
}
