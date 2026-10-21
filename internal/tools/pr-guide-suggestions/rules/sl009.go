// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl009 = &Rule{
	ID:       "SL009",
	Name:     "unit-in-naming",
	Severity: Warning,
	Fixable:  true,
	Check:    checkSL009,
}

// unitSuffixes are units of measure that should be written as "_in_<unit>".
// Singular/ambiguous forms (day, hour, min) are omitted to avoid false positives.
var unitSuffixes = []string{
	"bytes",
	"kib", "mib", "gib", "tib",
	"kb", "mb", "gb", "tb", "pb",
	"kbps", "mbps", "gbps",
	"ms", "seconds", "secs", "sec", "minutes", "hours", "hrs", "days",
}

func checkSL009(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		for _, unit := range unitSuffixes {
			if !strings.HasSuffix(n.Name, "_"+unit) {
				continue
			}
			if strings.HasSuffix(n.Name, "_in_"+unit) {
				break
			}
			prefix := strings.TrimSuffix(n.Name, "_"+unit)
			if prefix == "" {
				break
			}
			preferred := prefix + "_in_" + unit
			report(n,
				fmt.Sprintf("property %q uses a bare unit suffix; prefer the %q naming convention (%q)", n.Path, "_in_"+unit, preferred),
				&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
			)
			break
		}
	}
}
