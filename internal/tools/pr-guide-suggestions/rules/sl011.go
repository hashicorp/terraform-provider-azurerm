// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl011 = &Rule{
	ID:       "SL011",
	Name:     "redundant-is-prefix",
	Severity: Warning,
	Fixable:  true,
	Check:    checkSL011,
}

func checkSL011(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || n.Schema.ValueType() != typeBool {
			continue
		}
		if !strings.HasPrefix(n.Name, "is_") {
			continue
		}
		preferred := strings.TrimPrefix(n.Name, "is_")
		if preferred == "" {
			continue
		}

		report(n,
			fmt.Sprintf("boolean property %q has a redundant \"is_\" prefix (%q)", n.Path, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
	}
}
