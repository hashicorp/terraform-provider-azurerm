// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl006 = &Rule{
	ID:       "SL006",
	Name:     "block-needs-constraint",
	Severity: Warning,
	Check:    checkSL006,
}

func checkSL006(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || !userSettable(n.Schema) || len(n.Children) == 0 {
			continue
		}

		skip := false
		for _, c := range n.Children {
			// An opaque child (defined via a helper call) may itself be required
			// or constrained; skip to avoid a false positive.
			if c.Schema == nil {
				skip = true
				break
			}
			if c.Schema.Bool(fieldRequired) || c.Schema.Declares(fieldAtLeastOneOf) || c.Schema.Declares(fieldExactlyOneOf) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		report(n, fmt.Sprintf("block %q has no required fields and no AtLeastOneOf/ExactlyOneOf; add conditional validation so it cannot be configured empty", n.Path), nil)
	}
}
