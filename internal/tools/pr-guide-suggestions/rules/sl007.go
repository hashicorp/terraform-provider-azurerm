// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl007 = &Rule{
	ID:       "SL007",
	Name:     "array-limits",
	Severity: Warning,
	Check:    checkSL007,
}

func checkSL007(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || !userSettable(n.Schema) {
			continue
		}
		if !n.Schema.IsCollection() || !n.Schema.ElemIsScalarSchema() {
			continue
		}
		if setNonZeroInt(n.Schema, fieldMaxItems) {
			continue
		}

		report(n, fmt.Sprintf("array property %q does not set MaxItems; declare MinItems/MaxItems based on the API constraints", n.Path), nil)
	}
}
