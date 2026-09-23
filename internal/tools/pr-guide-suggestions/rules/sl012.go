// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl012 = &Rule{
	ID:       "SL012",
	Name:     "redundant-suffix",
	Severity: Warning,
	Fixable:  true,
	Check:    checkSL012,
}

var redundantSuffixes = []string{"_properties", "_config", "_profile"}

func checkSL012(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		for _, suffix := range redundantSuffixes {
			if !strings.HasSuffix(n.Name, suffix) {
				continue
			}
			preferred := strings.TrimSuffix(n.Name, suffix)
			if preferred == "" {
				break
			}
			report(n,
				fmt.Sprintf("property %q has a redundant %q suffix (%q)", n.Path, suffix, preferred),
				&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
			)
			break
		}
	}
}
