// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl008 = &Rule{
	ID:       "SL008",
	Name:     "sku-field-naming",
	Severity: Warning,
	Check:    checkSL008,
}

func checkSL008(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if !n.TopLevel || !strings.HasPrefix(n.Name, "sku_") {
			continue
		}

		count := 0
		for name := range n.Siblings {
			if strings.HasPrefix(name, "sku_") {
				count++
			}
		}
		if count <= 1 {
			continue
		}

		report(n, fmt.Sprintf("argument %q uses the sku_* format; with multiple sku_* arguments, prefer a single `sku` block", n.Path), nil)
	}
}
