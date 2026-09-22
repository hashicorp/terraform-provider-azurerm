// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl001 = &Rule{
	ID:       "SL001",
	Name:     "property-description-required",
	Severity: Warning,
	Check:    checkSL001,
}

func checkSL001(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil {
			continue
		}
		if n.Schema.Declares(fieldDescription) && !emptyStringLiteral(n.Schema.FieldValue(fieldDescription)) {
			continue
		}
		report(n, fmt.Sprintf("property %q is missing a description", n.Path), nil)
	}
}
