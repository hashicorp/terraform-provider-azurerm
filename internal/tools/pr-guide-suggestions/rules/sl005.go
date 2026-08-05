// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl005 = &Rule{
	ID:       "SL005",
	Name:     "validation-required",
	Severity: Warning,
	Check:    checkSL005,
}

func checkSL005(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || !userSettable(n.Schema) || hasValidation(n.Schema) {
			continue
		}

		switch n.Schema.ValueType() {
		case typeString:
			report(n, fmt.Sprintf("string property %q has no validation; add a ValidateFunc (StringIsNotEmpty at minimum)", n.Path), nil)
		case typeInt, typeFloat:
			report(n, fmt.Sprintf("numeric property %q has no validation; specify a valid range at minimum", n.Path), nil)
		}
	}
}
