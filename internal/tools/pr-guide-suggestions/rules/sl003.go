// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl003 = &Rule{
	ID:       "SL003",
	Name:     "limits-on-non-collection",
	Severity: Error,
	Fixable:  true,
	Check:    checkSL003,
}

func checkSL003(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || n.Schema.IsCollection() {
			continue
		}
		if !setNonZeroInt(n.Schema, fieldMinItems) && !setNonZeroInt(n.Schema, fieldMaxItems) {
			continue
		}

		report(n,
			fmt.Sprintf("property %q (%s) sets MinItems/MaxItems, which only apply to TypeList/TypeSet", n.Path, n.Schema.ValueType()),
			&Fix{Suggestion: fmt.Sprintf("remove MinItems/MaxItems from %q", n.Path)},
		)
	}
}
